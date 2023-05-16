package memory_test

import (
	"testing"

	"github.com/thiagolcmelo/webcrawler/src/events"
	"github.com/thiagolcmelo/webcrawler/src/memory"
)

func assertEventInMemoryEvents(
	t *testing.T,
	me *memory.Events,
	address string,
	eventType events.EventType,
	success bool,
	value int,
	expectedOccurrences int,
) {
	allEvents := me.GetReport()

	addressEvents, ok := allEvents[address]
	if !ok {
		t.Errorf("no events for address %s", address)
		return
	}

	actualOccurrences := 0
	for _, eventInstance := range addressEvents {
		if eventInstance.EventType == eventType && eventInstance.Success == success && eventInstance.Value == value {
			actualOccurrences++
		}
	}

	if expectedOccurrences != actualOccurrences {
		t.Errorf("event expected %d times, found %d times", expectedOccurrences, actualOccurrences)
	}
}

func TestMemoryEvents_LogDiscoveryEvent(t *testing.T) {
	me := memory.NewEvents()
	me.LogDiscoveryEvent("url1", true)
	me.LogDiscoveryEvent("url1", true)
	me.LogDiscoveryEvent("url2", true)
	me.LogDiscoveryEvent("url2", false)
	me.LogDiscoveryEvent("url3", false)
	me.LogDiscoveryEvent("url3", false)

	assertEventInMemoryEvents(t, me, "url1", events.Discovery, true, 0, 2)
	assertEventInMemoryEvents(t, me, "url1", events.Discovery, false, 0, 0)
	assertEventInMemoryEvents(t, me, "url2", events.Discovery, true, 0, 1)
	assertEventInMemoryEvents(t, me, "url2", events.Discovery, false, 0, 1)
	assertEventInMemoryEvents(t, me, "url3", events.Discovery, true, 0, 0)
	assertEventInMemoryEvents(t, me, "url3", events.Discovery, false, 0, 2)
}

func TestMemoryEvents_LogDownloadEvent(t *testing.T) {
	me := memory.NewEvents()
	me.LogDownloadEvent("url1", true)
	me.LogDownloadEvent("url1", true)
	me.LogDownloadEvent("url2", true)
	me.LogDownloadEvent("url2", false)
	me.LogDownloadEvent("url3", false)
	me.LogDownloadEvent("url3", false)

	assertEventInMemoryEvents(t, me, "url1", events.Download, true, 0, 2)
	assertEventInMemoryEvents(t, me, "url1", events.Download, false, 0, 0)
	assertEventInMemoryEvents(t, me, "url2", events.Download, true, 0, 1)
	assertEventInMemoryEvents(t, me, "url2", events.Download, false, 0, 1)
	assertEventInMemoryEvents(t, me, "url3", events.Download, true, 0, 0)
	assertEventInMemoryEvents(t, me, "url3", events.Download, false, 0, 2)
}

func TestMemoryEvents_LogParseEvent(t *testing.T) {
	me := memory.NewEvents()
	me.LogParseEvent("url1", true, 10)
	me.LogParseEvent("url1", true, 10)
	me.LogParseEvent("url2", true, 10)
	me.LogParseEvent("url2", false, 0)
	me.LogParseEvent("url3", false, 0)
	me.LogParseEvent("url3", false, 0)

	assertEventInMemoryEvents(t, me, "url1", events.Parse, true, 10, 2)
	assertEventInMemoryEvents(t, me, "url1", events.Parse, false, 0, 0)
	assertEventInMemoryEvents(t, me, "url2", events.Parse, true, 10, 1)
	assertEventInMemoryEvents(t, me, "url2", events.Parse, false, 0, 1)
	assertEventInMemoryEvents(t, me, "url3", events.Parse, true, 0, 0)
	assertEventInMemoryEvents(t, me, "url3", events.Parse, false, 0, 2)
}

func TestMemoryEvents_LogStoreEvent(t *testing.T) {
	me := memory.NewEvents()
	me.LogStoreEvent("url1", true)
	me.LogStoreEvent("url1", true)
	me.LogStoreEvent("url2", true)
	me.LogStoreEvent("url2", false)
	me.LogStoreEvent("url3", false)
	me.LogStoreEvent("url3", false)

	assertEventInMemoryEvents(t, me, "url1", events.Store, true, 0, 2)
	assertEventInMemoryEvents(t, me, "url1", events.Store, false, 0, 0)
	assertEventInMemoryEvents(t, me, "url2", events.Store, true, 0, 1)
	assertEventInMemoryEvents(t, me, "url2", events.Store, false, 0, 1)
	assertEventInMemoryEvents(t, me, "url3", events.Store, true, 0, 0)
	assertEventInMemoryEvents(t, me, "url3", events.Store, false, 0, 2)
}

func TestMemoryEvents_LogDispatchEvent(t *testing.T) {
	me := memory.NewEvents()
	me.LogDispatchEvent("url1", true, 5)
	me.LogDispatchEvent("url1", true, 5)
	me.LogDispatchEvent("url2", true, 5)
	me.LogDispatchEvent("url2", false, 0)
	me.LogDispatchEvent("url3", false, 0)
	me.LogDispatchEvent("url3", false, 0)

	assertEventInMemoryEvents(t, me, "url1", events.Dispatch, true, 5, 2)
	assertEventInMemoryEvents(t, me, "url1", events.Dispatch, false, 0, 0)
	assertEventInMemoryEvents(t, me, "url2", events.Dispatch, true, 5, 1)
	assertEventInMemoryEvents(t, me, "url2", events.Dispatch, false, 0, 1)
	assertEventInMemoryEvents(t, me, "url3", events.Dispatch, true, 5, 0)
	assertEventInMemoryEvents(t, me, "url3", events.Dispatch, false, 0, 2)
}

func TestMemoryEvents_GetReport(t *testing.T) {
	me := memory.NewEvents()
	me.LogDiscoveryEvent("url1", true)
	me.LogDownloadEvent("url1", true)
	me.LogParseEvent("url1", true, 10)
	me.LogStoreEvent("url1", true)
	me.LogDispatchEvent("url1", true, 5)

	assertEventInMemoryEvents(t, me, "url1", events.Discovery, true, 0, 1)
	assertEventInMemoryEvents(t, me, "url1", events.Download, true, 0, 1)
	assertEventInMemoryEvents(t, me, "url1", events.Parse, true, 10, 1)
	assertEventInMemoryEvents(t, me, "url1", events.Store, true, 0, 1)
	assertEventInMemoryEvents(t, me, "url1", events.Dispatch, true, 5, 1)
}

func TestMemoryEvents_IsAlreadyDiscovered(t *testing.T) {
	me := memory.NewEvents()
	me.LogDownloadEvent("url1", true)
	me.LogParseEvent("url1", true, 10)
	me.LogStoreEvent("url1", true)
	me.LogDispatchEvent("url1", true, 5)
	me.LogDiscoveryEvent("url2", true)

	type testCase struct {
		testName                    string
		me                          *memory.Events
		url                         string
		expectedIsAlreadyDiscovered bool
	}

	testCases := []testCase{
		{
			testName:                    "should_download_when_unknown",
			me:                          me,
			url:                         "new-url",
			expectedIsAlreadyDiscovered: true,
		},
		{
			testName:                    "should_download_when_not_discovered",
			me:                          me,
			url:                         "url1",
			expectedIsAlreadyDiscovered: true,
		},
		{
			testName:                    "should_not_download_when_discovered",
			me:                          me,
			url:                         "url2",
			expectedIsAlreadyDiscovered: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			actualIsAlreadyDiscovered := tc.me.IsAlreadyDiscovered(tc.url)
			if tc.expectedIsAlreadyDiscovered != actualIsAlreadyDiscovered {
				t.Errorf("expected %v, got %v", tc.expectedIsAlreadyDiscovered, actualIsAlreadyDiscovered)
			}
		})
	}
}
