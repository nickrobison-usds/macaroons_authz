package helpers

import (
	"encoding/base64"

	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/uuid"
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
