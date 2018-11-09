package helpers

import (
	"encoding/base64"
)

// BinaryToString converts a Byte array to a Base64 encoded string
func BinaryToString(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}
