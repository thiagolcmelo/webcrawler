package memory

import (
	"sync"
	"time"

	"github.com/thiagolcmelo/webcrawler/src/events"
)

type MemoryEvents struct {
	events map[string][]events.EventInstance
	sync.RWMutex
}

func NewMemoryEvents() *MemoryEvents {
	return &MemoryEvents{
		events: map[string][]events.EventInstance{},
	}
}

func (ms *MemoryEvents) addAddressIfNeeded(address string) {
	if _, ok := ms.events[address]; !ok {
		ms.events[address] = []events.EventInstance{}
	}
}

func (ms *MemoryEvents) LogDiscoveryEvent(address string, success bool) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Discovery,
		Success:   success,
		Time:      time.Now(),
	})
}

func (ms *MemoryEvents) LogDownloadEvent(address string, success bool) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Download,
		Success:   success,
		Time:      time.Now(),
	})
}

func (ms *MemoryEvents) LogParseEvent(address string, success bool, children int) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Parse,
		Success:   success,
		Value:     children,
		Time:      time.Now(),
	})
}

func (ms *MemoryEvents) LogStoreEvent(address string, success bool) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Store,
		Success:   success,
		Time:      time.Now(),
	})
}

func (ms *MemoryEvents) LogDispatchEvent(address string, success bool, children int) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Dispatch,
		Success:   success,
		Value:     children,
		Time:      time.Now(),
	})
}

func (ms *MemoryEvents) ShouldDownload(address string) bool {
	ms.Lock()
	defer ms.Unlock()

	addressEvents, ok := ms.events[address]
	if !ok {
		return true
	}
	for _, evt := range addressEvents {
		if evt.EventType == events.Discovery {
			return false
		}
	}
	return true
}

func (ms *MemoryEvents) GetReport() map[string][]events.EventInstance {
	ms.Lock()
	defer ms.Unlock()
	return ms.events
}
