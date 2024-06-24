package pebble

import (
	"github.com/cockroachdb/pebble"

	"github.com/gotd/td/session"

	"github.com/gotd/contrib/auth/kv"
)

var _ session.Storage = SessionStorage{}

// SessionStorage is a MTProto session Pebble storage.
type SessionStorage struct {
	kv.Session
}

// NewSessionStorage creates new SessionStorage.
func NewSessionStorage(db *pebble.DB, key string) SessionStorage {
	s := pebbleStorage{db: db, opts: pebble.Sync}
	return SessionStorage{
		Session: kv.NewSession(s, key),
	}
}
