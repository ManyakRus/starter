package bbolt

import (
	"go.etcd.io/bbolt"

	"github.com/gotd/td/session"

	"github.com/gotd/contrib/auth/kv"
)

var _ session.Storage = SessionStorage{}

// SessionStorage is a MTProto session bbolt storage.
type SessionStorage struct {
	kv.Session
}

// NewSessionStorage creates new SessionStorage.
func NewSessionStorage(db *bbolt.DB, key string, bucket []byte) SessionStorage {
	s := bboltStorage{db: db, bucket: bucket}
	return SessionStorage{
		Session: kv.NewSession(s, key),
	}
}
