package storage

import "github.com/thiagolcmelo/webcrawler/src/content"

// Storage defines an interface for storing downloaded content
type Storage interface {
	Add(content.Content) error
	UpdateContent(content.Content) error
	GetContent(string) (content.Content, error)
	GetAllContent() []content.Content
	IsRepeatedContent(content.Content) bool
}
