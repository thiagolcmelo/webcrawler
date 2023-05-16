package memory_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thiagolcmelo/webcrawler/src/content"
	"github.com/thiagolcmelo/webcrawler/src/memory"
)

func getMemoryStorageWithSamples(t *testing.T, samples map[string]string) *memory.Storage {
	sampleStorage := memory.NewStorage()

	for url, body := range samples {
		c, err := content.NewContentWithBody(url, []byte(body))
		if err != nil {
			t.Fatal(err)
		}

		err = sampleStorage.Add(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	if len(samples) != len(sampleStorage.GetAllContent()) {
		t.Fatal("error creating memory storage with samples")
	}

	return sampleStorage
}

func TestMemoryStorage_Add(t *testing.T) {
	type testCase struct {
		testName      string
		url           string
		body          string
		memoryStorage *memory.Storage
		expectedErr   error
	}

	samples := map[string]string{
		"http://url1.com": "content for url1",
		"http://url2.com": "content for url2",
		"http://url3.com": "content for url3",
	}

	testCases := []testCase{
		{
			testName:      "add_new_content_to_empty_storage_works",
			url:           "http://new-url.com",
			body:          "new content",
			memoryStorage: memory.NewStorage(),
			expectedErr:   nil,
		},
		{
			testName:      "add_new_content_works",
			url:           "http://new-url.com",
			body:          "new content",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   nil,
		},
		{
			testName:      "add_repeated_url_fails",
			url:           "http://url1.com",
			body:          "new content",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   memory.ErrAddingDuplicateURL,
		},
		{
			testName:      "add_repeated_content_fails",
			url:           "http://new-url.com",
			body:          "content for url1",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   memory.ErrAddingDuplicateContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			expectedContent, err := content.NewContentWithBody(tc.url, []byte(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			err = tc.memoryStorage.Add(expectedContent)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr != nil {
				return
			}

			actualContent, err := tc.memoryStorage.GetContent(expectedContent.Address)
			if err != nil {
				t.Errorf("content was not persisted")
			}

			if !reflect.DeepEqual(expectedContent, actualContent) {
				t.Errorf("content stored differs from original")
			}
		})
	}
}

func TestMemoryStorage_UpdateContent(t *testing.T) {
	type testCase struct {
		testName      string
		url           string
		body          string
		memoryStorage *memory.Storage
		expectedErr   error
	}

	samples := map[string]string{
		"http://url1.com": "content for url1",
		"http://url2.com": "content for url2",
		"http://url3.com": "content for url3",
	}

	testCases := []testCase{
		{
			testName:      "update_unknown_content_fails",
			url:           "http://new-url.com",
			body:          "new content",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   memory.ErrUpdatingUnknownContent,
		},
		{
			testName:      "update_existing_content_works",
			url:           "http://url1.com",
			body:          "new content for url1",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			expectedContent, err := content.NewContentWithBody(tc.url, []byte(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			err = tc.memoryStorage.UpdateContent(expectedContent)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr != nil {
				return
			}

			actualContent, err := tc.memoryStorage.GetContent(expectedContent.Address)
			if err != nil {
				t.Errorf("content was not persisted")
			}

			if !reflect.DeepEqual(expectedContent, actualContent) {
				t.Errorf("content stored differs from original")
			}
		})
	}
}

func TestMemoryStorage_GetContent(t *testing.T) {
	type testCase struct {
		testName      string
		url           string
		body          string
		memoryStorage *memory.Storage
		expectedErr   error
	}

	samples := map[string]string{
		"http://url1.com": "content for url1",
		"http://url2.com": "content for url2",
		"http://url3.com": "content for url3",
	}

	testCases := []testCase{
		{
			testName:      "get_unknown_content_fails",
			url:           "http://new-url.com",
			body:          "new content",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   memory.ErrGettingUnknownContent,
		},
		{
			testName:      "get_existing_content_fails",
			url:           "http://url1.com",
			body:          "content for url1",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			expectedContent, err := content.NewContentWithBody(tc.url, []byte(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			actualContent, err := tc.memoryStorage.GetContent(expectedContent.Address)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr != nil {
				return
			}

			if !reflect.DeepEqual(expectedContent, actualContent) {
				t.Errorf("content stored differs from original")
			}
		})
	}
}

func TestMemoryStorage_GetAllContent(t *testing.T) {
	type testCase struct {
		testName        string
		memoryStorage   *memory.Storage
		expectedContent map[string]string
	}

	samples := map[string]string{
		"http://url1.com": "content for url1",
		"http://url2.com": "content for url2",
		"http://url3.com": "content for url3",
	}

	testCases := []testCase{
		{
			testName:        "get_all_content_works_in_empty_storage",
			memoryStorage:   memory.NewStorage(),
			expectedContent: map[string]string{},
		},
		{
			testName:        "get_all_content_works_in_non_empty_storage",
			memoryStorage:   getMemoryStorageWithSamples(t, samples),
			expectedContent: samples,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			expectedContent := []content.Content{}
			for url, body := range tc.expectedContent {
				c, err := content.NewContentWithBody(url, []byte(body))
				if err != nil {
					t.Fatal(err)
				}
				expectedContent = append(expectedContent, c)
			}

			actualContent := tc.memoryStorage.GetAllContent()

			less := func(a, b content.Content) bool { return a.Address < b.Address }
			if diff := cmp.Diff(actualContent, expectedContent, cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("expected %#v, got %#v", actualContent, expectedContent)
			}
		})
	}
}

func TestMemoryStorage_IsRepeatedContent(t *testing.T) {
	type testCase struct {
		testName      string
		url           string
		body          string
		memoryStorage *memory.Storage
		expected      bool
	}

	samples := map[string]string{
		"http://url1.com": "content for url1",
		"http://url2.com": "content for url2",
		"http://url3.com": "content for url3",
	}

	testCases := []testCase{
		{
			testName:      "new_content_in_empty_storage_is_not_repeated",
			url:           "http://new-url.com",
			body:          "new content",
			memoryStorage: memory.NewStorage(),
			expected:      false,
		},
		{
			testName:      "new_content_in_non_empty_storage_is_not_repeated",
			url:           "http://new-url.com",
			body:          "new content",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expected:      false,
		},
		{
			testName:      "repeated_content_in_non_empty_storage_is_repeated_with_same_url",
			url:           "http://url1.com",
			body:          "content for url1",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expected:      true,
		},
		{
			testName:      "repeated_content_in_non_empty_storage_is_repeated_with_new_url",
			url:           "http://new-url.com",
			body:          "content for url1",
			memoryStorage: getMemoryStorageWithSamples(t, samples),
			expected:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			content, err := content.NewContentWithBody(tc.url, []byte(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			actual := tc.memoryStorage.IsRepeatedContent(content)

			if actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
