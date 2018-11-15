package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/envy"
	. "github.com/logrusorgru/aurora"
	macaroon "gopkg.in/macaroon.v2"
)

var acoID = "7680b5db-5cec-47a5-8c07-5dec827c3680"

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
		fmt.Println(string(cav.Id))
		fmt.Println(cav.Location)
	}

	// Try to make a request to read the data
	fmt.Println(Green("Trying to fetch the data"))

	var client http.Client
	url := fmt.Sprintf("http://localhost:8080/api/acos/test/%s?token=%s", acoID, token)
	fmt.Println(Blue(url))
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Status: %s. %s", resp.Status, body)

	/*
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
	*/
}
