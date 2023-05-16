package memory

import (
	"errors"
	"sync"

	"github.com/thiagolcmelo/webcrawler/src/content"
	"golang.org/x/exp/maps"
)

var (
	// ErrAddingDuplicateURL should be used when trying to add a reapeated url
	ErrAddingDuplicateURL = errors.New("could not add content because url exists already")
	// ErrAddingDuplicateContent should be used when trying to add reapeated content
	ErrAddingDuplicateContent = errors.New("could not add content because it exists already")
	// ErrUpdatingUnknownContent should be used when trying to update unknown content
	ErrUpdatingUnknownContent = errors.New("could not update content because it does not exist yet")
	// ErrGettingUnknownContent should be used when trying to get unknown content
	ErrGettingUnknownContent = errors.New("could not get content because it does not exist yet")
)

// Storage is an in memory implementation of Storage
type Storage struct {
	existingChecksums map[[32]byte]struct{}
	urlToChecksum     map[string][32]byte
	cache             map[string]content.Content
	sync.RWMutex
}

// NewStorage is a factory for in memory Storage
func NewStorage() *Storage {
	return &Storage{
		existingChecksums: map[[32]byte]struct{}{},
		urlToChecksum:     map[string][32]byte{},
		cache:             map[string]content.Content{},
	}
}

// Add adds a new content to storage
func (ms *Storage) Add(c content.Content) error {
	ms.Lock()
	defer ms.Unlock()
	if _, ok := ms.cache[c.Address]; ok {
		return ErrAddingDuplicateURL
	}
	if _, ok := ms.existingChecksums[c.BodyHash]; ok {
		return ErrAddingDuplicateContent
	}
	ms.cache[c.Address] = c
	ms.existingChecksums[c.BodyHash] = struct{}{}

	return nil
}

// UpdateContent updates an existing content
func (ms *Storage) UpdateContent(c content.Content) error {
	ms.Lock()
	defer ms.Unlock()
	if _, ok := ms.cache[c.Address]; !ok {
		return ErrUpdatingUnknownContent
	}
	ms.cache[c.Address] = c
	return nil
}

// GetContent retrieves a content by its address
func (ms *Storage) GetContent(address string) (content.Content, error) {
	ms.Lock()
	defer ms.Unlock()
	c, ok := ms.cache[address]
	if !ok {
		return content.Content{}, ErrGettingUnknownContent
	}
	return c, nil
}

// GetAllContent returns a list with all existing content
func (ms *Storage) GetAllContent() []content.Content {
	ms.Lock()
	defer ms.Unlock()
	return maps.Values(ms.cache)
}

// IsRepeatedContent checks if there is a content with the same URL or same body in memory
func (ms *Storage) IsRepeatedContent(c content.Content) bool {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.existingChecksums[c.BodyHash]; !ok {
		return false
	}

	// log for future use
	ms.urlToChecksum[c.Address] = c.BodyHash

	return true
}
