package bbolt

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-faster/errors"
	"go.etcd.io/bbolt"

	"github.com/gotd/contrib/storage"
)

var _ storage.PeerStorage = PeerStorage{}

// PeerStorage is a peer storage based on pebble.
type PeerStorage struct {
	bbolt  *bbolt.DB
	bucket []byte
}

// NewPeerStorage creates new peer storage using bbolt.
func NewPeerStorage(db *bbolt.DB, bucket []byte) *PeerStorage {
	return &PeerStorage{bbolt: db, bucket: bucket}
}

type bboltIterator struct {
	tx      *bbolt.Tx
	iter    *bbolt.Cursor
	lastErr error
	value   storage.Peer
}

func (p *bboltIterator) Close() error {
	return p.tx.Rollback()
}

func (p *bboltIterator) Next(ctx context.Context) bool {
Next:
	k, v := p.iter.Next()
	if v == nil {
		return false
	}

	for {
		if bytes.HasPrefix(k, storage.PeerKeyPrefix) {
			break
		}

		k, v = p.iter.Next()
		if v == nil {
			return false
		}
	}

	if err := json.Unmarshal(v, &p.value); err != nil {
		if errors.Is(err, storage.ErrPeerUnmarshalMustInvalidate) {
			goto Next // skip
		}
		p.lastErr = errors.Wrap(err, "unmarshal")
		return false
	}

	return true
}

func (p *bboltIterator) Err() error {
	return p.lastErr
}

func (p *bboltIterator) Value() storage.Peer {
	return p.value
}

// Iterate creates and returns new PeerIterator.
func (s PeerStorage) Iterate(ctx context.Context) (storage.PeerIterator, error) {
	tx, err := s.bbolt.Begin(false)
	if err != nil {
		return nil, errors.Errorf("create tx: %w", err)
	}

	bucket := tx.Bucket(s.bucket)
	if bucket == nil {
		return nil, errors.Errorf("bucket %q does not exist", s.bucket)
	}

	cur := bucket.Cursor()
	cur.Seek(storage.PeerKeyPrefix)
	cur.Prev()
	return &bboltIterator{
		tx:   tx,
		iter: cur,
	}, nil
}

func (s PeerStorage) add(associated []string, value storage.Peer) (err error) {
	err = s.bbolt.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(s.bucket)
		if err != nil {
			return errors.Errorf("create bucket: %w", err)
		}

		data, err := json.Marshal(value)
		if err != nil {
			return errors.Errorf("marshal: %w", err)
		}
		id := storage.KeyFromPeer(value).Bytes(nil)

		if err := bucket.Put(id, data); err != nil {
			return errors.Errorf("set id <-> data: %w", err)
		}

		for _, key := range associated {
			if err := bucket.Put([]byte(key), id); err != nil {
				return errors.Errorf("set key <-> id: %w", err)
			}
		}

		return nil
	})
	return
}

// Add adds given peer to the storage.
func (s PeerStorage) Add(ctx context.Context, value storage.Peer) error {
	return s.add(value.Keys(), value)
}

// Find finds peer using given key.
func (s PeerStorage) Find(ctx context.Context, key storage.PeerKey) (p storage.Peer, rerr error) {
	rerr = s.bbolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(s.bucket)
		if bucket == nil {
			return errors.Errorf("bucket %q does not exist", s.bucket)
		}

		data := bucket.Get(key.Bytes(nil))
		if data == nil {
			return storage.ErrPeerNotFound
		}

		if err := json.Unmarshal(data, &p); err != nil {
			if errors.Is(err, storage.ErrPeerUnmarshalMustInvalidate) {
				return storage.ErrPeerNotFound
			}
			return errors.Errorf("unmarshal: %w", err)
		}
		return nil
	})
	return
}

// Assign adds given peer to the storage and associate it to the given key.
func (s PeerStorage) Assign(ctx context.Context, key string, value storage.Peer) error {
	return s.add(append(value.Keys(), key), value)
}

// Resolve finds peer using associated key.
func (s PeerStorage) Resolve(ctx context.Context, key string) (p storage.Peer, rerr error) {
	rerr = s.bbolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(s.bucket)
		if bucket == nil {
			return errors.Errorf("bucket %q does not exist", s.bucket)
		}

		id := bucket.Get([]byte(key))
		if id == nil {
			return storage.ErrPeerNotFound
		}

		data := bucket.Get(id)
		if data == nil {
			return storage.ErrPeerNotFound
		}

		if err := json.Unmarshal(data, &p); err != nil {
			if errors.Is(err, storage.ErrPeerUnmarshalMustInvalidate) {
				return storage.ErrPeerNotFound
			}
			return errors.Errorf("unmarshal: %w", err)
		}
		return nil
	})
	return
}
