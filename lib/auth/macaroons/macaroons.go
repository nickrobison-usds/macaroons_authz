package macaroons

import (
	"encoding/base64"

	"github.com/gobuffalo/pop/nulls"
	"gopkg.in/macaroon-bakery.v2/bakery"
	macaroon "gopkg.in/macaroon.v2"
)

var oven *bakery.Oven

// DischargeRequest holds the basic structure of a caveat discharge request
type DischargeRequest struct {
	Id        string `httprequest:"id,form,omitempty"`
	Id64      string `httprequest:"id64,form,omitempty"`
	Caveat    string `httprequest:"caveat64,form,omitempty"`
	Token     string `httprequest:"token,form,omitempty"`
	Token64   string `httprequest:"token64,form,omitempty"`
	TokenKind string `httprequest:"token-kind,form,omitempty"`
}

func init() {
	oven = bakery.NewOven(bakery.OvenParams{})
}

// DecodeMacaroon returns a bakery.Macaroon from a given string
// The string can either be raw JSON, or a base64 encoded string
func DecodeMacaroon(s string) (*bakery.Macaroon, error) {
	// Check to see if we're JSON, or not
	if s[0] == '{' {
		return MacaroonFromJSON(s)
	}
	b, err := macaroon.Base64Decode([]byte(s))
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

// MacaroonFromJSON returns a bakery.Macaroon from a JSON string
func MacaroonFromJSON(s string) (*bakery.Macaroon, error) {
	var m macaroon.Macaroon
	err := m.UnmarshalJSON([]byte(s))
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
