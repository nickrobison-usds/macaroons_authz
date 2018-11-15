package macaroons

import (
	"encoding/base64"

	"github.com/gobuffalo/pop/nulls"
	"gopkg.in/macaroon-bakery.v2/bakery"
	macaroon "gopkg.in/macaroon.v2"
)

var oven *bakery.Oven

func init() {
	oven = bakery.NewOven(bakery.OvenParams{})
}

// DecodeMacaroon returns a bakery.Macaroon from a base64 encoded string
func DecodeMacaroon(s string) (*bakery.Macaroon, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return MacaroonFromBytes(b)
}

// MacaroonFromBytes returns a bakery.Macaroon from a byte array
func MacaroonFromBytes(b []byte) (*bakery.Macaroon, error) {
	var m macaroon.Macaroon
	err := m.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}

	return bakery.NewLegacyMacaroon(&m)
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
