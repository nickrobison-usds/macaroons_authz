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

var acoID = "8752e6a2-2838-4b8b-91d8-9642e8dae729"
var userID = "280ea269-c427-48e3-bd97-4aa255fbf610"

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
