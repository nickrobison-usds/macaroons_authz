package macaroons

import (
	"context"

	"gopkg.in/macaroon-bakery.v2/bakery"
)

var key = []byte("test key")

// DevRootKeyStore is a simple way of testing functionality. It always returns a test key.
type DevRootKeyStore struct {
}

func (s *DevRootKeyStore) Get(_ context.Context, id []byte) ([]byte, error) {
	return key, nil
}

func (s *DevRootKeyStore) RootKey(_ context.Context) (rootKey, id []byte, err error) {
	return key, []byte("0"), nil
}

func NewDevKeyRootStore() bakery.RootKeyStore {
	return new(DevRootKeyStore)
}
