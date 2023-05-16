package basic

import (
	"github.com/thiagolcmelo/webcrawler/src/events"
	"github.com/thiagolcmelo/webcrawler/src/frontier"
)

// Dispatcher is a basic implementation for Dispatcher interface
type Dispatcher struct {
	events   events.Events
	frontier frontier.Frontier
}

// NewDispatcher creates a new dispatcher
func NewDispatcher(events events.Events, frontier frontier.Frontier) *Dispatcher {
	return &Dispatcher{
		events:   events,
		frontier: frontier,
	}
}

// DispatchNewUrls dispatches new URLs to the download frontier
func (bd *Dispatcher) DispatchNewUrls(urls []string) (int, error) {
	newUrls := []string{}

	for _, url := range urls {
		if bd.events.IsAlreadyDiscovered(url) {
			newUrls = append(newUrls, url)
		}
	}

	for _, url := range newUrls {
		bd.frontier.Publish(url)
	}

	return len(newUrls), nil
}
