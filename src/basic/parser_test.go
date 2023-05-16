package basic_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thiagolcmelo/webcrawler/src/basic"
	"github.com/thiagolcmelo/webcrawler/src/content"
)

func TestParser_Parse(t *testing.T) {
	type testCase struct {
		testName      string
		url           string
		body          string
		expectedLinks []string
		expectedErr   error
	}

	testCases := []testCase{
		{
			testName:      "empty_content_returns_nothing",
			url:           "http://domain.com",
			body:          "",
			expectedLinks: []string{},
			expectedErr:   nil,
		},
		{
			testName:      "find_full_links",
			url:           "http://domain.com",
			body:          "<a href=\"http://domain.com/path\">link</a>",
			expectedLinks: []string{"http://domain.com/path"},
			expectedErr:   nil,
		},
		{
			testName:      "complete_omitted_scheme_to_same_domain",
			url:           "http://domain.com",
			body:          "<a href=\"//domain.com/path\">link</a>",
			expectedLinks: []string{"http://domain.com/path"},
			expectedErr:   nil,
		},
		{
			testName:      "ignore_omitted_scheme_to_other_domain",
			url:           "http://domain.com",
			body:          "<a href=\"//other.com/path\">link</a>",
			expectedLinks: []string{},
			expectedErr:   nil,
		},
		{
			testName:      "find_relative_links",
			url:           "http://domain.com",
			body:          "<a href=\"/path\">link</a>",
			expectedLinks: []string{"http://domain.com/path"},
			expectedErr:   nil,
		},
		{
			testName:      "ignore_other_domain",
			url:           "http://domain.com",
			body:          "<a href=\"http://otherdomain.com/path\">link</a>",
			expectedLinks: []string{},
			expectedErr:   nil,
		},
		{
			testName:      "find_valids_and_ignore_other_domain",
			url:           "http://domain.com",
			body:          "<a href=\"http://otherdomain.com/path\">link</a><a href=\"/path1\">link</a><a href=\"http://domain.com/path2\">link</a>",
			expectedLinks: []string{"http://domain.com/path1", "http://domain.com/path2"},
			expectedErr:   nil,
		},
		{
			testName:      "dedupe_repeated_absolutes",
			url:           "http://domain.com",
			body:          "<a href=\"http://domain.com/path1\">link</a><a href=\"http://domain.com/path1\">link</a>",
			expectedLinks: []string{"http://domain.com/path1"},
			expectedErr:   nil,
		},
		{
			testName:      "dedupe_repeated_relatives",
			url:           "http://domain.com",
			body:          "<a href=\"/path1\">link</a><a href=\"/path1\">link</a>",
			expectedLinks: []string{"http://domain.com/path1"},
			expectedErr:   nil,
		},
		{
			testName:      "dedupe_repeated_absolute_and_relative",
			url:           "http://domain.com",
			body:          "<a href=\"/path1\">link</a><a href=\"http://domain.com/path1\">link</a>",
			expectedLinks: []string{"http://domain.com/path1"},
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			c, err := content.NewContent(tc.url)
			if err != nil {
				t.Fatal(err)
			}
			c.Body = []byte(tc.body)

			parser := basic.NewParser()

			err = parser.Parse(&c)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			less := func(a, b string) bool { return a < b }
			if diff := cmp.Diff(tc.expectedLinks, c.GetChildrenList(), cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("expected %#v, got %#v", tc.expectedLinks, c.GetChildrenList())
			}
		})
	}
}
