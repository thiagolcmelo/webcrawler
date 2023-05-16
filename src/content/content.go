package content

import (
	"crypto/sha256"
	"net/url"

	"golang.org/x/exp/maps"
)

// Content bundles a URL, its info, and also the content associated with it
type Content struct {
	Address     string
	Body        []byte
	BodyHash    [32]byte
	Children    map[string]struct{}
	ContentType string
	*url.URL
}

// NewContent generates a content object for a URL
func NewContent(address string) (Content, error) {
	url, err := url.Parse(address)
	if err != nil {
		return Content{}, err
	}

	if url.Path == "" {
		url.Path = "/"
	}
	url.Fragment = ""

	return Content{
		url.String(),
		[]byte{},
		[32]byte{},
		map[string]struct{}{},
		"",
		url,
	}, nil
}

// NewContentWithBody generates a content object for a URL with its known body
func NewContentWithBody(address string, body []byte) (Content, error) {
	c, err := NewContent(address)
	if err != nil {
		return c, err
	}
	c.Body = body
	c.CreateChecksum()
	return c, nil
}

// GetChildrenList returns a list with all relevant links existing in the body of the associated content
func (c Content) GetChildrenList() []string {
	return maps.Keys(c.Children)
}

// CreateChecksum creates a checksum for the body content
func (c *Content) CreateChecksum() {
	c.BodyHash = sha256.Sum256(c.Body)
}
