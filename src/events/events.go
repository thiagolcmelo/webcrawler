package events

import (
	"fmt"
	"time"
)

type EventType int

const (
	Discovery EventType = iota
	Download
	Parse
	Store
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

type EventInstance struct {
	EventType EventType
	Success   bool
	Value     int
	Time      time.Time
}

type Events interface {
	LogDiscoveryEvent(string, bool)
	LogDownloadEvent(string, bool)
	LogParseEvent(string, bool, int)
	LogStoreEvent(string, bool)
	LogDispatchEvent(string, bool, int)
	GetReport() map[string][]EventInstance
	ShouldDownload(string) bool
}
