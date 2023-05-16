package downloader

import (
	"context"

	"github.com/thiagolcmelo/webcrawler/src/content"
)

// Downloader defines an interface for downloading content
type Downloader interface {
	Download(context.Context, *content.Content) error
}
