package macaroons

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gobuffalo/uuid"
	macaroon "gopkg.in/macaroon.v1"
)

func MacaroonFromBytes(b []byte) (macaroon.Macaroon, error) {
	var m macaroon.Macaroon
	err := m.UnmarshalBinary(b)
	return m, err
}

// DelegateACOToUser restricts an existing macaroon to a certain user
func DelegateACOToUser(acoID uuid.UUID, userID uuid.UUID, m *macaroon.Macaroon) (*macaroon.Macaroon, error) {

	nonce, err := GenerateNonce()
	if err != nil {
		return m, err
	}

	locString := fmt.Sprintf("http://localhost:8080/api/aco/%s/verify/%s", acoID.String(), userID.String())
	err = m.AddThirdPartyCaveat([]byte("test key"), nonce, locString)
	if err != nil {
		return m, err
	}

	// Add another test caveat
	err = m.AddThirdPartyCaveat([]byte("test key"), "0test nonce", "http://localhost:8080/other/things/to/test")
	if err != nil {
		return nil, err
	}

	return m, nil
}

// GenerateNonce creates a random ID that can be used for macaroons.
func GenerateNonce() (string, error) {
	nonce := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(nonce), err
}

func EncodeMacaroon(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}
