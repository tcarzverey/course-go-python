package storage

import (
	"context"
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

// Store is the storage interface used by all handler code.
// Implementations: MemoryStorage (steps 0–7), PGXStorage (step8).
type Store interface {
	Save(ctx context.Context, originalURL string) (string, error)
	Get(ctx context.Context, code string) (*URLRecord, error)
	IncrementClicks(ctx context.Context, code string) error
}

type MemoryStorage struct {
	mu   sync.RWMutex
	urls map[string]*URLRecord
}

func New() Store {
	return &MemoryStorage{
		urls: make(map[string]*URLRecord),
	}
}

func (s *MemoryStorage) Save(_ context.Context, originalURL string) (string, error) {
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

func (s *MemoryStorage) Get(_ context.Context, code string) (*URLRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.urls[code]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStorage) IncrementClicks(_ context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.urls[code]; ok {
		r.Clicks++
	}
	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
