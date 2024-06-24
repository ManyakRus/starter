package pebble

import (
	"github.com/cockroachdb/pebble"

	"github.com/gotd/contrib/auth/kv"
)

// Credentials stores user credentials to Pebble.
type Credentials struct {
	kv.Credentials
}

// NewCredentials creates new Credentials.
func NewCredentials(db *pebble.DB) Credentials {
	s := pebbleStorage{db: db, opts: pebble.Sync}
	return Credentials{
		Credentials: kv.NewCredentials(s),
	}
}
