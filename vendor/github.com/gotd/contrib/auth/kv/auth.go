package kv

import (
	"context"

	"github.com/go-faster/errors"

	tgauth "github.com/gotd/td/telegram/auth"
)

// Credentials is a generic implementation of credential storage
// over key-value Storage.
type Credentials struct {
	storage               Storage
	phoneKey, passwordKey string
}

// NewCredentials creates new Credentials.
func NewCredentials(storage Storage) Credentials {
	return Credentials{
		storage:     storage,
		phoneKey:    "phone",
		passwordKey: "password",
	}
}

// WithPhoneKey sets phone key to use.
func (c Credentials) WithPhoneKey(phoneKey string) Credentials {
	c.phoneKey = phoneKey
	return c
}

// WithPasswordKey sets password key to use.
func (c Credentials) WithPasswordKey(passwordKey string) Credentials {
	c.passwordKey = passwordKey
	return c
}

// Phone implements Credentials and returns phone.
func (c Credentials) Phone(ctx context.Context) (string, error) {
	return c.storage.Get(ctx, c.phoneKey)
}

// Password implements Credentials and returns password.
func (c Credentials) Password(ctx context.Context) (string, error) {
	r, err := c.storage.Get(ctx, c.passwordKey)
	if errors.Is(err, ErrKeyNotFound) {
		return r, tgauth.ErrPasswordNotProvided
	}
	return r, err
}

// SavePhone stores given phone to storage.
func (c Credentials) SavePhone(ctx context.Context, phone string) error {
	return c.storage.Set(ctx, c.phoneKey, phone)
}

// SavePassword stores given password to storage.
func (c Credentials) SavePassword(ctx context.Context, password string) error {
	return c.storage.Set(ctx, c.passwordKey, password)
}
