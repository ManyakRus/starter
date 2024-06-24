package storage

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/telegram/query/dialogs"
	"github.com/gotd/td/tg"
)

// FindPeer finds peer using given storage.
func FindPeer(ctx context.Context, s PeerStorage, p tg.PeerClass) (Peer, error) {
	var key dialogs.DialogKey

	if err := key.FromPeer(p); err != nil {
		return Peer{}, err
	}

	return s.Find(ctx, PeerKey{
		Kind: key.Kind,
		ID:   key.ID,
	})
}

// ForEach calls callback on every iterator element.
func ForEach(ctx context.Context, iterator PeerIterator, cb func(Peer) error) error {
	for iterator.Next(ctx) {
		if err := cb(iterator.Value()); err != nil {
			return errors.Errorf("callback: %w", err)
		}
	}
	return iterator.Err()
}
