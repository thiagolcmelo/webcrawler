package src_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thiagolcmelo/webcrawler/src"
	"github.com/thiagolcmelo/webcrawler/src/memory"
)

type webpage struct {
	url              string
	body             string
	expectedChildren []string
}

func sampleWebsite(url string) []webpage {
	webpage0 := webpage{
		url: fmt.Sprintf("%s/", url),
		body: fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Page 1 Title</title>
		</head>
		<body>
			<a href="/page1">Relative link to Page 1</a>
			<a href="%s/page1">Absolute link to Page 1</a>
			<a href="/page2">Relative link to Page 2</a>
			<a href="%s/page2">Absolute link to Page 2</a>
			<a href="/page3">Relative link to Page 3</a>
			<a href="%s/page3">Absolute link to Page 3</a>
			<a href="http://link-to-somewhere-else.com/">Link to somewhere else</a>
		</body>
		</html>`, url, url, url),
		expectedChildren: []string{fmt.Sprintf("%s/page1", url), fmt.Sprintf("%s/page2", url), fmt.Sprintf("%s/page3", url)},
	}

	webpage1 := webpage{
		url: fmt.Sprintf("%s/page1", url),
		body: fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Page 1 Title</title>
		</head>
		<body>
			<a href="/">Relative link to Page 0</a>
			<a href="%s/">Absolute link to Page 0</a>
			<a href="/page2">Relative link to Page 2</a>
			<a href="%s/page2">Absolute link to Page 2</a>
			<a href="/page3">Relative link to Page 3</a>
			<a href="%s/page3">Absolute link to Page 3</a>
			<a href="http://link-to-somewhere-else.com/">Link to somewhere else</a>
		</body>
		</html>`, url, url, url),
		expectedChildren: []string{fmt.Sprintf("%s/", url), fmt.Sprintf("%s/page2", url), fmt.Sprintf("%s/page3", url)},
	}

	webpage2 := webpage{
		url: fmt.Sprintf("%s/page2", url),
		body: fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Page 2 Title</title>
		</head>
		<body>
			<a href="/page1">Relative link to Page 1</a>
			<a href="%s/page1">Absolute link to Page 1</a>
			<a href="http://link-to-somewhere-else.com/">Link to somewhere else</a>
		</body>
		</html>`, url),
		expectedChildren: []string{fmt.Sprintf("%s/page1", url)},
	}

	webpage3 := webpage{
		url: fmt.Sprintf("%s/page3", url),
		body: `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Page 3 Title</title>
		</head>
		<body>
			<a href="http://link-to-somewhere-else.com/">Link to somewhere else</a>
		</body>
		</html>`,
		expectedChildren: []string{},
	}

	return []webpage{webpage0, webpage1, webpage2, webpage3}
}

func TestOrchestrator(t *testing.T) {
	website := map[string]webpage{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		url := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
		page, ok := website[url]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page.body))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	for _, page := range sampleWebsite(server.URL) {
		website[page.url] = page
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	var buf bytes.Buffer

	frontier := memory.NewFrontier()
	storage := memory.NewStorage()
	events := memory.NewEvents()

	orchestrator := src.NewOrchestrator(ctx, 10, frontier, storage, events, 1, time.Second, 2)
	orchestrator.Start(server.URL)
	err := orchestrator.PrintReport(&buf, true, false)
	if err != nil {
		t.Fatal(err)
	}

	var actualResult []src.OrchestratorOutputItem
	err = json.Unmarshal(buf.Bytes(), &actualResult)
	if err != nil {
		t.Fatal(err)
	}

	if len(actualResult) != len(website) {
		t.Errorf("expected %d results, got %d", len(website), len(actualResult))
	}

	less := func(a, b string) bool { return a < b }
	for _, resultItem := range actualResult {
		if page, ok := website[resultItem.URL]; !ok {
			t.Errorf("unknown result %s", resultItem.URL)
		} else {
			if diff := cmp.Diff(resultItem.Children, page.expectedChildren, cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("expected %#v, got %#v", page.expectedChildren, resultItem.Children)
			}
		}
	}
}
