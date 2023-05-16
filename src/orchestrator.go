package src

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/thiagolcmelo/webcrawler/src/basic"
	"github.com/thiagolcmelo/webcrawler/src/content"
	"github.com/thiagolcmelo/webcrawler/src/dispatcher"
	"github.com/thiagolcmelo/webcrawler/src/downloader"
	"github.com/thiagolcmelo/webcrawler/src/events"
	"github.com/thiagolcmelo/webcrawler/src/frontier"
	"github.com/thiagolcmelo/webcrawler/src/parser"
	"github.com/thiagolcmelo/webcrawler/src/storage"
)

type OrchestratorOutputItem struct {
	Url         string   `json:"url"`
	ContentType string   `json:"contentType"`
	Children    []string `json:"children"`
}

type Orchestrator struct {
	ctx         context.Context
	wg          sync.WaitGroup
	downloaders int
	frontier    frontier.Frontier
	storage     storage.Storage
	events      events.Events
	downloader  downloader.Downloader
	parser      parser.Parser
	dispatcher  dispatcher.Dispatcher
}

func NewOrchestrator(
	ctx context.Context,
	downloaders int,
	frontier frontier.Frontier,
	storage storage.Storage,
	events events.Events,
	retries int,
	backoff time.Duration,
	backoffMultiplier int,
) *Orchestrator {
	return &Orchestrator{
		ctx:         ctx,
		wg:          sync.WaitGroup{},
		downloaders: downloaders,
		frontier:    frontier,
		storage:     storage,
		events:      events,
		downloader:  basic.NewBasicDownloader(retries, backoff, backoffMultiplier),
		parser:      basic.NewBasicParser(),
		dispatcher:  basic.NewBasicDispatcher(events, frontier),
	}
}

func (o *Orchestrator) Start(seed string) {
	for i := 0; i < o.downloaders; i++ {
		go func(i int) {
			for {
				select {
				case url := <-o.frontier.Consume():
					go o.processUrl(url)
				case <-o.ctx.Done():
					return
				}
			}
		}(i)
	}
	// wg is decremented when processUrl finishes
	o.wg.Add(1)
	o.frontier.Publish(seed)

	// wait for downloads complete or context to be canceled
	c := make(chan struct{})
	go func() {
		defer close(c)
		o.wg.Wait()
	}()
	select {
	case <-c:
		log.Printf("all urls were processed")
	case <-o.ctx.Done():
		log.Printf("timeout exceeded")
	}
}

func (o *Orchestrator) discovery(c *content.Content) error {
	if c.Scheme == "" {
		o.wg.Add(1)
		go o.processUrl(fmt.Sprintf("https://%s", c.Address))
		o.wg.Add(1)
		go o.processUrl(fmt.Sprintf("http://%s", c.Address))

		o.events.LogDiscoveryEvent(c.Address, false)
		return fmt.Errorf("url missing schema [%s], trying https and http", c.Address)
	}

	if !o.events.ShouldDownload(c.Address) {
		o.events.LogDiscoveryEvent(c.Address, false)
		return fmt.Errorf("repeated url [%s]", c.Address)
	}

	o.events.LogDiscoveryEvent(c.Address, true)
	return nil
}

func (o *Orchestrator) download(c *content.Content) error {
	err := o.downloader.Download(o.ctx, c)
	if err != nil {
		o.events.LogDownloadEvent(c.Address, false)
		return fmt.Errorf("download failed: %v", err)
	}
	o.events.LogDownloadEvent(c.Address, true)
	return nil
}

func (o *Orchestrator) skipRepeated(c *content.Content) error {
	if o.storage.IsRepeatedContent(*c) {
		return fmt.Errorf("repeated content for url [%s]", c.Address)
	}
	return nil
}

func (o *Orchestrator) parse(c *content.Content) error {
	err := o.parser.Parse(c)
	if err != nil {
		o.events.LogParseEvent(c.Address, false, 0)
		return fmt.Errorf("parse failed: %v", err)
	}
	o.events.LogParseEvent(c.Address, true, len(c.Children))
	return nil
}

func (o *Orchestrator) store(c *content.Content) error {
	err := o.storage.Add(*c)
	if err != nil {
		o.events.LogStoreEvent(c.Address, false)
		return fmt.Errorf("store failed: %v", err)
	}
	o.events.LogStoreEvent(c.Address, true)
	return nil
}

func (o *Orchestrator) dispatch(c *content.Content) error {
	n, err := o.dispatcher.DispatchNewUrls(c.GetChildrenList())
	if err != nil {
		o.events.LogDispatchEvent(c.Address, false, n)
		return fmt.Errorf("dispatch failed: %v", err)
	}
	o.events.LogDispatchEvent(c.Address, true, n)
	// wg is decremented when processUrl finishes
	o.wg.Add(n)
	return nil
}

func (o *Orchestrator) processUrl(url string) {
	// wg is incremented upon adding seed and after dispatching
	defer o.wg.Done()

	c, err := content.NewContent(url)
	if err != nil {
		log.Printf("error parsing url [%s]: %v", url, err)
	}

	type action func(*content.Content) error

	var actions []action = []action{
		o.discovery,
		o.download,
		o.skipRepeated,
		o.parse,
		o.store,
		o.dispatch,
	}

	for _, a := range actions {
		err := a(&c)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (o *Orchestrator) PrintReport(w io.Writer, isJson bool, isIndented bool) error {
	allContent := o.storage.GetAllContent()
	sort.Sort(sortByAddress(allContent))

	output := make([]OrchestratorOutputItem, len(allContent))
	for i, c := range allContent {
		output[i] = OrchestratorOutputItem{
			Url:         c.Address,
			ContentType: c.ContentType,
			Children:    c.GetChildrenList(),
		}
	}

	if isJson {
		var jsonData []byte
		var err error
		if isIndented {
			jsonData, err = json.MarshalIndent(output, "", "    ")
		} else {
			jsonData, err = json.Marshal(output)
		}
		if err != nil {
			return err
		}
		_, err = w.Write(jsonData)
		return err
	}

	for _, content := range allContent {
		if _, err := w.Write([]byte(fmt.Sprintf("%s\n", content.Address))); err != nil {
			return err
		}
		for _, child := range content.GetChildrenList() {
			if _, err := w.Write([]byte(fmt.Sprintf("  |- %s\n", child))); err != nil {
				return err
			}
		}
	}

	return nil
}

type sortByAddress []content.Content

func (a sortByAddress) Len() int           { return len(a) }
func (a sortByAddress) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByAddress) Less(i, j int) bool { return a[i].Address < a[j].Address }
