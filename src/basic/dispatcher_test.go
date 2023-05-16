package basic_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thiagolcmelo/webcrawler/src/basic"
	"github.com/thiagolcmelo/webcrawler/src/events"
	"golang.org/x/exp/maps"
)

type fakeFrontier struct {
	Items []string
}

func (ff *fakeFrontier) Publish(url string) error {
	ff.Items = append(ff.Items, url)
	return nil
}

func (ff *fakeFrontier) Consume() <-chan string {
	return make(<-chan string)
}

type fakeEvents struct {
	UrlsToDownload map[string]struct{}
}

func (fe *fakeEvents) LogDiscoveryEvent(string, bool)     {}
func (fe *fakeEvents) LogDownloadEvent(string, bool)      {}
func (fe *fakeEvents) LogParseEvent(string, bool, int)    {}
func (fe *fakeEvents) LogStoreEvent(string, bool)         {}
func (fe *fakeEvents) LogDispatchEvent(string, bool, int) {}
func (fe *fakeEvents) GetReport() map[string][]events.EventInstance {
	return map[string][]events.EventInstance{}
}
func (fe *fakeEvents) IsAlreadyDiscovered(url string) bool {
	_, ok := fe.UrlsToDownload[url]
	return ok
}

func TestBasicDispatcher_DispatchNewUrls(t *testing.T) {
	type testCase struct {
		testName       string
		shouldDownload map[string]bool
		expectedErr    error
		expectedUrls   []string
	}

	testCases := []testCase{
		{
			testName:       "dispatch_nothing_from_empty",
			shouldDownload: map[string]bool{},
			expectedErr:    nil,
			expectedUrls:   []string{},
		},
		{
			testName:       "dispatch_expected",
			shouldDownload: map[string]bool{"url1": true, "url2": false, "url3": true},
			expectedErr:    nil,
			expectedUrls:   []string{"url1", "url3"},
		},
		{
			testName:       "dispatch_nothing_if_nothing_should_be_dispatched",
			shouldDownload: map[string]bool{"url1": false, "url2": false, "url3": false},
			expectedErr:    nil,
			expectedUrls:   []string{},
		},
		{
			testName:       "dispatch_all_if_all_should_be_dispatched",
			shouldDownload: map[string]bool{"url1": true, "url2": true, "url3": true},
			expectedErr:    nil,
			expectedUrls:   []string{"url1", "url2", "url3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ff := &fakeFrontier{
				Items: []string{},
			}
			fe := &fakeEvents{
				UrlsToDownload: map[string]struct{}{},
			}
			for url, shouldDownload := range tc.shouldDownload {
				if shouldDownload {
					fe.UrlsToDownload[url] = struct{}{}
				}
			}

			dispatcher := basic.NewDispatcher(fe, ff)
			_, err := dispatcher.DispatchNewUrls(maps.Keys(tc.shouldDownload))
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			actualUrls := ff.Items

			less := func(a, b string) bool { return a < b }
			if diff := cmp.Diff(actualUrls, tc.expectedUrls, cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("expected %#v, got %#v", actualUrls, tc.expectedUrls)
			}
		})
	}
}
