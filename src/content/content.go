package content

import (
	"crypto/sha256"
	"net/url"

	"golang.org/x/exp/maps"
)

type Content struct {
	Address     string
	Body        []byte
	BodyHash    [32]byte
	Children    map[string]struct{}
	ContentType string
	*url.URL
}

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

func NewContentWithBody(address string, body []byte) (Content, error) {
	c, err := NewContent(address)
	if err != nil {
		return c, err
	}
	c.Body = body
	c.CreateChecksum()
	return c, nil
}

func (c Content) GetChildrenList() []string {
	return maps.Keys(c.Children)
}

func (c *Content) CreateChecksum() {
	c.BodyHash = sha256.Sum256(c.Body)
}
