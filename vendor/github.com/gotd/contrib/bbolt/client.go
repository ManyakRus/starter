package bbolt

import (
	"context"

	"github.com/go-faster/errors"
	"go.etcd.io/bbolt"

	"github.com/gotd/contrib/auth/kv"
)

type bboltStorage struct {
	db     *bbolt.DB
	bucket []byte
}

func (p bboltStorage) Set(ctx context.Context, k, v string) (rerr error) {
	return p.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(p.bucket)
		if err != nil {
			return errors.Errorf("create bucket: %w", err)
		}

		if err := bucket.Put([]byte(k), []byte(v)); err != nil {
			return errors.Errorf("put: %w", err)
		}
		return nil
	})
}

func (p bboltStorage) Get(ctx context.Context, k string) (r string, err error) {
	err = p.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(p.bucket)
		if bucket == nil {
			return errors.Errorf("bucket %q does not exist", p.bucket)
		}

		result := bucket.Get([]byte(k))
		if result == nil {
			return kv.ErrKeyNotFound
		}

		r = string(result)
		return nil
	})
	return
}
