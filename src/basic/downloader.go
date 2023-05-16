package basic

import (
	"context"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/thiagolcmelo/webcrawler/src/content"
)

var (
	// ErrExecutingRequest should be used when a request fails
	ErrExecutingRequest = errors.New("could not execute request")
	// ErrResponseStatusNotOK should be used when the status code is not ok
	ErrResponseStatusNotOK = errors.New("response status not 200")
)

// Downloader is a basic implementation of the Downloader interface
type Downloader struct {
	retries           int
	backoff           time.Duration
	backoffMultiplier int
}

// NewDownloader is a factory for basic.Downloader
func NewDownloader(retries int, backoff time.Duration, backoffMultiplier int) *Downloader {
	return &Downloader{
		retries:           retries,
		backoff:           backoff,
		backoffMultiplier: backoffMultiplier,
	}
}

// Download attempts to fetch a URL content and store in the provided content object
func (bd *Downloader) Download(ctx context.Context, c *content.Content) error {
	for i := 0; i < bd.retries; i++ {
		err := bd.download(ctx, c)
		// if there was an error and it is not the last attempt
		if err != nil && i < bd.retries-1 {
			log.Printf("attempt %d for url [%s] failed due to %v", i+1, c.Address, err)

			multiplier := math.Pow(float64(bd.backoffMultiplier), float64(i))

			time.Sleep(bd.backoff * time.Duration(multiplier))
			continue
		}
		return err
	}
	return nil
}

func (bd *Downloader) download(ctx context.Context, c *content.Content) error {
	// create a request that can be canceled from the context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Address, nil)
	if err != nil {
		return err
	}

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			return ErrExecutingRequest
		}
		return nil
	}
	defer resp.Body.Close()

	// check if the response is 200
	if resp.StatusCode != http.StatusOK {
		return ErrResponseStatusNotOK
	}

	// store the downloaded response in the content body
	c.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.CreateChecksum()

	// store the content type in the content as well
	c.ContentType = resp.Header.Get("Content-Type")

	return nil
}
