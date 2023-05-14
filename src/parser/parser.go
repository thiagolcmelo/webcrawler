package parser

import "github.com/thiagolcmelo/webcrawler/src/content"

type Parser interface {
	Parse(*content.Content) error
}
