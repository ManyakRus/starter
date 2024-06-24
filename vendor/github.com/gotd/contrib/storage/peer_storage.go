package storage

import (
	"context"
	"io"

	"github.com/go-faster/errors"
)

// ErrPeerNotFound is a special error to return when peer not found.
var ErrPeerNotFound = errors.New("peer not found")

// PeerStorage is abstraction for peer storage.
type PeerStorage interface {
	// Add adds given peer to the storage.
	Add(ctx context.Context, value Peer) error
	// Find finds peer using given key.
	// If peer not found, it returns ErrPeerNotFound error.
	Find(ctx context.Context, key PeerKey) (Peer, error)

	// Assign adds given peer to the storage and associates it to the given key.
	Assign(ctx context.Context, key string, value Peer) error
	// Resolve finds peer using associated key.
	// If peer not found, it returns ErrPeerNotFound error.
	Resolve(ctx context.Context, key string) (Peer, error)

	// Iterate creates and returns new PeerIterator.
	Iterate(ctx context.Context) (PeerIterator, error)
}

// PeerIterator is a peer iterator.
type PeerIterator interface {
	Next(ctx context.Context) bool
	Err() error
	Value() Peer
	io.Closer
}
