package basic

import (
	"github.com/thiagolcmelo/webcrawler/src/events"
	"github.com/thiagolcmelo/webcrawler/src/frontier"
)

type BasicDispatcher struct {
	events   events.Events
	frontier frontier.Frontier
}

func NewBasicDispatcher(events events.Events, frontier frontier.Frontier) *BasicDispatcher {
	return &BasicDispatcher{
		events:   events,
		frontier: frontier,
	}
}

func (bd *BasicDispatcher) DispatchNewUrls(urls []string) (int, error) {
	newUrls := []string{}

	for _, url := range urls {
		if bd.events.ShouldDownload(url) {
			newUrls = append(newUrls, url)
		}
	}

	for _, url := range newUrls {
		bd.frontier.Publish(url)
	}

	return len(newUrls), nil
}
