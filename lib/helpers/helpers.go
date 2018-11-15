package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/uuid"
)

const (
	// NonceSize denotes the number of bytes used for the nonce
	NonceSize = 12
)

var log logger.FieldLogger

func init() {
	log = logger.NewLogger("HELPERS")
}

// BinaryToString converts a Byte array to a Base64 encoded string
func BinaryToString(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

// MustGenerateID forces the generation of a V4 UUID, otherwise it panics.
func MustGenerateID() uuid.UUID {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return u
}

// UUIDOfString always returns a uuid.UUID, otherwise it panics
func UUIDOfString(id string) uuid.UUID {
	str, err := uuid.FromString(id)
	if err != nil {
		panic(err)
	}
	return str
}

// GenerateNonce creates a random ID that can be used for macaroons.
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nonce, err
	}

	return nonce, err
}
