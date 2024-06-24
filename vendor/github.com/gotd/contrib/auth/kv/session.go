package kv

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/session"
)

var _ session.Storage = Session{}

// Session is a generic implementation of session storage
// over key-value Storage.
type Session struct {
	storage Storage
	key     string
}

// NewSession creates new Session.
func NewSession(storage Storage, key string) Session {
	return Session{storage: storage, key: key}
}

// LoadSession loads session using given key from storage.
func (s Session) LoadSession(ctx context.Context) ([]byte, error) {
	r, err := s.storage.Get(ctx, s.key)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return nil, session.ErrNotFound
		}
		return nil, err
	}

	return []byte(r), nil
}

// StoreSession saves session using given key to storage.
func (s Session) StoreSession(ctx context.Context, data []byte) error {
	return s.storage.Set(ctx, s.key, string(data))
}
