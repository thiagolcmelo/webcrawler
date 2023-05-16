package parser

import "github.com/thiagolcmelo/webcrawler/src/content"

// Parser defines an interface for parsing downloaded content
type Parser interface {
	Parse(*content.Content) error
}
