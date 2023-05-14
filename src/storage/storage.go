package storage

import "github.com/thiagolcmelo/webcrawler/src/content"

type Storage interface {
	Add(content.Content) error
	UpdateContent(content.Content) error
	GetContent(string) (content.Content, error)
	GetAllContent() []content.Content
	IsRepeatedContent(content.Content) bool
}
