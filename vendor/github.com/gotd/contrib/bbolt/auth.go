package bbolt

import (
	"go.etcd.io/bbolt"

	"github.com/gotd/contrib/auth/kv"
)

// Credentials stores user credentials to bbolt.
type Credentials struct {
	kv.Credentials
}

// NewCredentials creates new Credentials.
func NewCredentials(db *bbolt.DB, bucket []byte) Credentials {
	s := bboltStorage{db: db, bucket: bucket}
	return Credentials{
		Credentials: kv.NewCredentials(s),
	}
}
