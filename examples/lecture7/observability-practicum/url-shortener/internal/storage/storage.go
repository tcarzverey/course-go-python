package storage

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrNotFound = errors.New("URL not found")

type URLRecord struct {
	Code        string
	OriginalURL string
	CreatedAt   time.Time
	Clicks      int64
}

type Storage struct {
	mu   sync.RWMutex
	urls map[string]*URLRecord
}

func New() *Storage {
	return &Storage{
		urls: make(map[string]*URLRecord),
	}
}

func (s *Storage) Save(originalURL string) (string, error) {
	code := generateCode()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[code] = &URLRecord{
		Code:        code,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}
	return code, nil
}

func (s *Storage) Get(code string) (*URLRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.urls[code]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *Storage) IncrementClicks(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.urls[code]; ok {
		r.Clicks++
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
