package main

import (
	"encoding/base64"
	"fmt"
	macaroon "gopkg.in/macaroon.v2"
)

func main() {
	fmt.Println("Hello!")

	//Create a test Macaroon

	mac, err := macaroon.New([]byte("test-key"), []byte("test-id"), "http://test.loc", macaroon.V2)
	if err != nil {
		panic(err)
	}

	fmt.Println(serializeToString(mac))

	// Add a caveat

	err = mac.AddThirdPartyCaveat([]byte("SECRET for 3rd party caveat"), []byte("test-third-party"), "http://auth.mybank")
	if err != nil {
		panic(err)
	}

	fmt.Println(serializeToString(mac))

}

func serializeToString(m *macaroon.Macaroon) string {
	marsh, err := m.MarshalJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println("JSON: ", string(marsh))

	return base64.RawURLEncoding.EncodeToString(marsh)
}
