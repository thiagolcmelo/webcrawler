package content_test

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thiagolcmelo/webcrawler/src/content"
)

func TestContent_NewContent(t *testing.T) {
	valid, err := content.NewContent("http://complete-url.com/path")
	if err != nil {
		t.Fatal(err)
	}

	if valid.Address != "http://complete-url.com/path" {
		t.Errorf("invalid address %s", valid.Address)
	}
}

func TestContent_NewContentEnforceRoot(t *testing.T) {
	valid, err := content.NewContent("http://valid-url.com")
	if err != nil {
		t.Fatal(err)
	}

	if valid.Address != "http://valid-url.com/" {
		t.Errorf("invalid address %s", valid.Address)
	}
}

func TestContent_NewContentRemoveFragment(t *testing.T) {
	valid, err := content.NewContent("http://valid-url.com#somethingUseless")
	if err != nil {
		t.Fatal(err)
	}

	if valid.Address != "http://valid-url.com/" {
		t.Errorf("invalid address %s", valid.Address)
	}
}

func TestContent_NewContentWithBody(t *testing.T) {
	type testCase struct {
		testName         string
		url              string
		body             []byte
		expectedChecksum [32]byte
	}

	emptyBody := []byte("")
	emptyBodyChecksum := sha256.Sum256(emptyBody)

	nonEmptyBody := []byte("some valid body")
	nonEmptyBodyChecksum := sha256.Sum256(nonEmptyBody)

	testCases := []testCase{
		{
			testName:         "valid_url_empty_body",
			url:              "http://valid-url.com/",
			body:             emptyBody,
			expectedChecksum: emptyBodyChecksum,
		},
		{
			testName:         "valid_url_non_empty_body",
			url:              "http://valid-url.com/",
			body:             nonEmptyBody,
			expectedChecksum: nonEmptyBodyChecksum,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			c, err := content.NewContentWithBody(tc.url, tc.body)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(c.Body, tc.body) {
				t.Fatalf("wrong body, expected %s, got %s", string(tc.body), string(c.Body))
			}

			if c.BodyHash != tc.expectedChecksum {
				t.Fatalf("wrong checksum, expected %#v, got %#vs", tc.expectedChecksum, c.BodyHash)
			}
		})
	}
}

func TestContent_GetChildrenList(t *testing.T) {
	c, err := content.NewContent("http://valid-url.com")
	if err != nil {
		t.Fatal(err)
	}

	less := func(a, b string) bool { return a < b }

	if diff := cmp.Diff(c.GetChildrenList(), []string{}, cmpopts.SortSlices(less)); diff != "" {
		t.Errorf("expected %#v, got %#v", []string{}, c.GetChildrenList())
	}

	c.Children["child1"] = struct{}{}
	c.Children["child2"] = struct{}{}
	c.Children["child3"] = struct{}{}
	expectedChildren := []string{"child1", "child2", "child3"}

	if diff := cmp.Diff(c.GetChildrenList(), expectedChildren, cmpopts.SortSlices(less)); diff != "" {
		t.Errorf("expected %#v, got %#v", expectedChildren, c.GetChildrenList())
	}
}

func TestContent_CreateChecksum(t *testing.T) {
	type testCase struct {
		testName         string
		url              string
		body             []byte
		expectedChecksum [32]byte
	}

	emptyBody := []byte("")
	emptyBodyChecksum := sha256.Sum256(emptyBody)

	nonEmptyBody := []byte("some valid body")
	nonEmptyBodyChecksum := sha256.Sum256(nonEmptyBody)

	testCases := []testCase{
		{
			testName:         "valid_url_empty_body",
			url:              "http://valid-url.com/",
			body:             emptyBody,
			expectedChecksum: emptyBodyChecksum,
		},
		{
			testName:         "valid_url_non_empty_body",
			url:              "http://valid-url.com/",
			body:             nonEmptyBody,
			expectedChecksum: nonEmptyBodyChecksum,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			c, err := content.NewContent(tc.url)
			if err != nil {
				t.Fatal(err)
			}
			c.Body = tc.body
			c.CreateChecksum()

			if c.BodyHash != tc.expectedChecksum {
				t.Fatalf("wrong checksum, expected %#v, got %#vs", tc.expectedChecksum, c.BodyHash)
			}
		})
	}
}
