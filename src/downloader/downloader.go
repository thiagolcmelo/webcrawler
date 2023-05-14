package downloader

import (
	"context"

	"github.com/thiagolcmelo/webcrawler/src/content"
)

type Downloader interface {
	Download(context.Context, *content.Content) error
}
