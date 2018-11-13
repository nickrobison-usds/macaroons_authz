package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/envy"
	. "github.com/logrusorgru/aurora"
	macaroon "gopkg.in/macaroon.v2"
)

var acoID = "496290ec-3f8e-481c-b637-bc6f29a005bf"
var userID = "c4654147-7f16-43a2-a5eb-aeee4a74984a"

func main() {
	token, err := envy.MustGet("TOKEN")
	if err != nil {
		panic(err)
	}
	fmt.Println(Green("Starting up"))
	fmt.Println(Sprintf("Token: %s", Cyan(token)))

	var m macaroon.Macaroon
	bin, err := macaroon.Base64Decode([]byte(token))
	if err != nil {
		panic(err)
	}
	err = m.UnmarshalBinary(bin)
	if err != nil {
		panic(err)
	}

	cavs := m.Caveats()

	fmt.Println(Blue("Printing Caveats:"))

	for _, cav := range cavs {

		fmt.Println(cav.Location)
	}

	// lWe need to get an authorization discharge macaroon
	fmt.Println(Green("Fetching Authorization macaroon."))

	data := map[string]interface{}{
		"aco_id":   acoID,
		"user_id":  userID,
		"macaroon": bin,
	}

	jsonValues, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	var client http.Client
	_, err = client.Post("http://localhost:8080/api/acos/verify", "application/json", bytes.NewBuffer(jsonValues))
	if err != nil {
		panic(err)
	}
}
