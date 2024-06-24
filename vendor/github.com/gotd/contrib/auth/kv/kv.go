package kv

import (
	"context"

	"github.com/go-faster/errors"
)

// Storage represents generic KV storage.
type Storage interface {
	Set(ctx context.Context, k, v string) error
	Get(ctx context.Context, k string) (string, error)
}

// ErrKeyNotFound is a special error to return when given key not found.
var ErrKeyNotFound = errors.New("key not found")
