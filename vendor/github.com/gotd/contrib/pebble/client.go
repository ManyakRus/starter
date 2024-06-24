package pebble

import (
	"context"

	"github.com/cockroachdb/pebble"
	"github.com/go-faster/errors"
	"go.uber.org/multierr"

	"github.com/gotd/contrib/auth/kv"
)

type pebbleStorage struct {
	db   *pebble.DB
	opts *pebble.WriteOptions
}

func (p pebbleStorage) Set(ctx context.Context, k, v string) (rerr error) {
	b := p.db.NewBatch()
	defer func() {
		multierr.AppendInto(&rerr, b.Close())
	}()

	d := b.SetDeferred(len(k), len(v))
	copy(d.Key, k)
	copy(d.Value, v)
	_ = d.Finish()

	return b.Commit(p.opts)
}

func (p pebbleStorage) Get(ctx context.Context, k string) (string, error) {
	r, closer, err := p.db.Get([]byte(k))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return "", kv.ErrKeyNotFound
		}
		return "", err
	}
	v := string(r)

	return v, closer.Close()
}
