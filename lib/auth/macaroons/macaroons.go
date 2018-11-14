package macaroons

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gofrs/uuid"
	"gopkg.in/macaroon-bakery.v2/bakery"
	macaroon "gopkg.in/macaroon.v2"
)

const (
	// NonceSize denotes the number of bytes used for the nonce
	NonceSize = 12
)

var oven *bakery.Oven

func init() {
	oven = bakery.NewOven(bakery.OvenParams{})
}

func MacaroonFromBytes(b []byte) (*bakery.Macaroon, error) {
	var m macaroon.Macaroon
	err := m.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}

	return bakery.NewLegacyMacaroon(&m)
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

// DelegateACOToUser restricts an existing macaroon to a certain user
func DelegateACOToUser(acoID uuid.UUID, userID uuid.UUID, m *bakery.Macaroon) (*bakery.Macaroon, error) {

	// Add the User restrictions
	m2, err := AddThirdPartyCaveat(m, "http://localhost:8080/api/aco/verify", []string{fmt.Sprintf("user_id = %s", userID.String())})
	if err != nil {
		return nil, err
	}

	return m2, nil
}

func EncodeMacaroon(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func MacaroonToByteSlice(m *bakery.Macaroon) (nulls.ByteSlice, error) {
	b, err := m.M().MarshalBinary()
	if err != nil {
		return nulls.NewByteSlice([]byte{}), err
	}

	return nulls.NewByteSlice(b), nil
}
