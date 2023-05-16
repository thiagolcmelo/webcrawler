package events

import (
	"fmt"
	"time"
)

// EventType is an enumeration for defining types of events
type EventType int

const (
	// Discovery is used for a URL that was discovered
	Discovery EventType = iota
	// Download is used for a URL that was downloade
	Download
	// Parse is used for a URL that was parsed
	Parse
	// Store is used for a URL that was stored
	Store
	// Dispatch is used for a URL that was dispatched
	Dispatch
)

func (et EventType) String() string {
	switch et {
	case Discovery:
		return "discovery"
	case Download:
		return "download"
	case Parse:
		return "parse"
	case Store:
		return "store"
	case Dispatch:
		return "dispatch"
	default:
		return fmt.Sprintf("%d", int(et))
	}
}

// EventInstance bundles information about an event
type EventInstance struct {
	EventType EventType
	Success   bool
	Value     int
	Time      time.Time
}

// Events defines an interface for a storage used for logging events
type Events interface {
	LogDiscoveryEvent(string, bool)
	LogDownloadEvent(string, bool)
	LogParseEvent(string, bool, int)
	LogStoreEvent(string, bool)
	LogDispatchEvent(string, bool, int)
	GetReport() map[string][]EventInstance
	IsAlreadyDiscovered(string) bool
}
