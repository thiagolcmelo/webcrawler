package basic

import (
	"bytes"

	"github.com/thiagolcmelo/webcrawler/src/content"
	"golang.org/x/net/html"
)

// Parser is a basic implementation of the Parser interface
type Parser struct{}

// NewParser is a factory for basic.Parser`
func NewParser() *Parser {
	return &Parser{}
}

func (ep Parser) extractLinksFromData(data []byte) ([]string, error) {
	links := []string{}
	reader := bytes.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken, html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

// Parse updates a content object with the links existing in its body
func (ep *Parser) Parse(c *content.Content) error {
	links, err := ep.extractLinksFromData(c.Body)
	if err != nil {
		return err
	}

	// add scheme and hostname if needed
	for i, l := 0, len(links); i < l; i++ {

		linkAsContent, err := content.NewContent(links[i])
		if err != nil {
			continue
		}
		if linkAsContent.Hostname() == "" {
			linkAsContent.Host = c.Host
		} else if linkAsContent.Hostname() != c.Hostname() {
			continue
		}

		if linkAsContent.Scheme == "" {
			linkAsContent.Scheme = c.Scheme
		}

		c.Children[linkAsContent.String()] = struct{}{}
	}

	return nil
}
