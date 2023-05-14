package memory

import (
	"errors"
	"sync"

	"github.com/thiagolcmelo/webcrawler/src/content"
	"golang.org/x/exp/maps"
)

var (
	ErrAddingDuplicateUrl     = errors.New("could not add content because url exists already")
	ErrAddingDuplicateContent = errors.New("could not add content because it exists already")
	ErrUpdatingUnknownContent = errors.New("could not update content because it does not exist yet")
	ErrGettingUnknownContent  = errors.New("could not get content because it does not exist yet")
)

type MemoryStorage struct {
	existingChecksums map[[32]byte]struct{}
	urlToChecksum     map[string][32]byte
	cache             map[string]content.Content
	sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		existingChecksums: map[[32]byte]struct{}{},
		urlToChecksum:     map[string][32]byte{},
		cache:             map[string]content.Content{},
	}
}

func (ms *MemoryStorage) Add(c content.Content) error {
	ms.Lock()
	defer ms.Unlock()
	if _, ok := ms.cache[c.Address]; ok {
		return ErrAddingDuplicateUrl
	}
	if _, ok := ms.existingChecksums[c.BodyHash]; ok {
		return ErrAddingDuplicateContent
	}
	ms.cache[c.Address] = c
	ms.existingChecksums[c.BodyHash] = struct{}{}

	return nil
}

func (ms *MemoryStorage) UpdateContent(c content.Content) error {
	ms.Lock()
	defer ms.Unlock()
	if _, ok := ms.cache[c.Address]; !ok {
		return ErrUpdatingUnknownContent
	}
	ms.cache[c.Address] = c
	return nil
}

func (ms *MemoryStorage) GetContent(address string) (content.Content, error) {
	ms.Lock()
	defer ms.Unlock()
	c, ok := ms.cache[address]
	if !ok {
		return content.Content{}, ErrGettingUnknownContent
	}
	return c, nil
}

func (ms *MemoryStorage) GetAllContent() []content.Content {
	ms.Lock()
	defer ms.Unlock()
	return maps.Values(ms.cache)
}

func (ms *MemoryStorage) IsRepeatedContent(c content.Content) bool {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.existingChecksums[c.BodyHash]; !ok {
		return false
	}

	// log for future use
	ms.urlToChecksum[c.Address] = c.BodyHash

	return true
}
