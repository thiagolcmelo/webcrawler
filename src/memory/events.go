package memory

import (
	"sync"
	"time"

	"github.com/thiagolcmelo/webcrawler/src/events"
)

// Events is an in memory implementation of Events
type Events struct {
	events map[string][]events.EventInstance
	sync.RWMutex
}

// NewEvents is a factory for in memory Events
func NewEvents() *Events {
	return &Events{
		events: map[string][]events.EventInstance{},
	}
}

func (ms *Events) addAddressIfNeeded(address string) {
	if _, ok := ms.events[address]; !ok {
		ms.events[address] = []events.EventInstance{}
	}
}

// LogDiscoveryEvent adds a Discovery event to memory
func (ms *Events) LogDiscoveryEvent(address string, success bool) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Discovery,
		Success:   success,
		Time:      time.Now(),
	})
}

// LogDownloadEvent adds a Download event to memory
func (ms *Events) LogDownloadEvent(address string, success bool) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Download,
		Success:   success,
		Time:      time.Now(),
	})
}

// LogParseEvent adds a Parse event to memory
func (ms *Events) LogParseEvent(address string, success bool, children int) {
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

// LogStoreEvent adds a Store event to memory
func (ms *Events) LogStoreEvent(address string, success bool) {
	ms.Lock()
	defer ms.Unlock()
	ms.addAddressIfNeeded(address)
	ms.events[address] = append(ms.events[address], events.EventInstance{
		EventType: events.Store,
		Success:   success,
		Time:      time.Now(),
	})
}

// LogDispatchEvent adds a Dispatch event to memory
func (ms *Events) LogDispatchEvent(address string, success bool, children int) {
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

// IsAlreadyDiscovered informs if an address was already discovered
func (ms *Events) IsAlreadyDiscovered(address string) bool {
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

// GetReport simply return all events
func (ms *Events) GetReport() map[string][]events.EventInstance {
	ms.Lock()
	defer ms.Unlock()
	return ms.events
}
