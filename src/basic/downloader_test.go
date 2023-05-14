package basic_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/thiagolcmelo/webcrawler/src/basic"
	"github.com/thiagolcmelo/webcrawler/src/content"
)

func TestDownloader_Download(t *testing.T) {
	homeBody := []byte("hello from home")
	homeContentType := "text/html; charset=utf-8"
	homeBodyHash := sha256.Sum256(homeBody)

	contactBody := []byte("hello from contact")
	contactContentType := ""
	contactBodyHash := sha256.Sum256(contactBody)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.RequestURI == "/" {
			w.Header().Set("Content-Type", homeContentType)
			w.Write(homeBody)

		} else if r.RequestURI == "/contact" {
			w.Header().Set("Content-Type", contactContentType)
			w.Write(contactBody)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	type testCase struct {
		testName           string
		url                string
		expectedErr        error
		expectedBody       []byte
		expectedBodyHash   [32]byte
		expectedCntentType string
	}

	testCases := []testCase{
		{
			testName:           "request_to_valid_root_url_updates_content",
			url:                server.URL,
			expectedErr:        nil,
			expectedBody:       homeBody,
			expectedBodyHash:   homeBodyHash,
			expectedCntentType: homeContentType,
		},
		{
			testName:           "request_to_valid_url_updates_content",
			url:                fmt.Sprintf("%s/contact", server.URL),
			expectedErr:        nil,
			expectedBody:       contactBody,
			expectedBodyHash:   contactBodyHash,
			expectedCntentType: contactContentType,
		},
		{
			testName:           "request_to_bad_url_fails",
			url:                "tcp://localhost.com:",
			expectedErr:        basic.ErrExecutingRequest,
			expectedBody:       []byte{},
			expectedBodyHash:   [32]byte{},
			expectedCntentType: "",
		},
		{
			testName:           "bad_request_fails",
			url:                fmt.Sprintf("%s/invalid", server.URL),
			expectedErr:        basic.ErrResponseStatusNotOK,
			expectedBody:       []byte{},
			expectedBodyHash:   [32]byte{},
			expectedCntentType: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			c, err := content.NewContent(tc.url)
			if err != nil {
				t.Fatal(err)
			}

			downloader := basic.NewBasicDownloader(1, time.Second, 2)

			err = downloader.Download(context.Background(), &c)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if !bytes.Equal(c.Body, tc.expectedBody) {
				t.Errorf("body is not correct")
			}

			if c.BodyHash != tc.expectedBodyHash {
				t.Errorf("body hash is not correct")
			}

			if c.ContentType != tc.expectedCntentType {
				t.Errorf("content type is not correct")
			}
		})
	}
}

func TestDownloader_Retries(t *testing.T) {

	type testCase struct {
		testName       string
		retries        int
		successRequest int // if greater than retries, simulates failure
	}

	testCases := []testCase{
		{
			testName:       "works_in_first_request",
			retries:        1,
			successRequest: 1,
		},
		{
			testName:       "works_in_second_request",
			retries:        2,
			successRequest: 2,
		},
		{
			testName:       "stop_after_max_retries",
			retries:        2,
			successRequest: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				if attempts == tc.successRequest {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte("hello from home"))
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			c, err := content.NewContent(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			downloader := basic.NewBasicDownloader(tc.retries, 200*time.Millisecond, 2)

			err = downloader.Download(context.Background(), &c)
			if err != nil && attempts >= tc.successRequest {
				t.Errorf("it should be successful, got %v", err)
			}

			if attempts != tc.retries {
				t.Errorf("expected %d retries, it happened %d", tc.retries, attempts)
			}
		})
	}
}

func TestDownloader_Backoff(t *testing.T) {
	type testCase struct {
		testName          string
		retries           int
		successRequest    int // if greater than retries, simulates failure
		backoff           time.Duration
		backoffMultiplier int
	}

	testCases := []testCase{
		{
			testName:          "respects_backoff_and_succeed",
			retries:           3,
			successRequest:    3,
			backoff:           200 * time.Millisecond,
			backoffMultiplier: 2,
		},
		{
			testName:          "respects_backoff_and_fail",
			retries:           3,
			successRequest:    4,
			backoff:           200 * time.Millisecond,
			backoffMultiplier: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			attempts := 0
			currentBackoff := 1
			nextRetry := time.Now()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if time.Now().Before(nextRetry) {
					t.Errorf("retry was too soon")
				}

				attempts++
				if attempts == tc.successRequest {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte("hello from home"))
					w.WriteHeader(http.StatusOK)
					return
				}

				nextRetry = time.Now().Add(tc.backoff * time.Duration(currentBackoff))
				currentBackoff *= tc.backoffMultiplier

				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			c, err := content.NewContent(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			downloader := basic.NewBasicDownloader(tc.retries, 200*time.Millisecond, 2)

			err = downloader.Download(context.Background(), &c)
			if err != nil && attempts >= tc.successRequest {
				t.Errorf("it should be successful, got %v", err)
			}

			if attempts != tc.retries {
				t.Errorf("expected %d retries, it happened %d", tc.retries, attempts)
			}
		})
	}
}
