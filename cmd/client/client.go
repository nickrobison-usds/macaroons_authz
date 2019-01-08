package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/logrusorgru/aurora"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/macaroons"
	"gopkg.in/macaroon-bakery.v2/httpbakery"
)

var acoID = "4d24935c-a5bb-4ab1-a4ba-0ac7e730e4a9"
var userID = "97952122-60cb-4d04-8d00-a3fbca8b53e4"

func main() {
	/*
		token, err := envy.MustGet("TOKEN")
		if err != nil {
			panic(err)
		}
		fmt.Println(Green("Starting up"))

		// Decode it?

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
	*/

	// Get the token

	req1, err := http.NewRequest("GET", "http://localhost:3002/token", nil)
	if err != nil {
		panic(err)
	}

	q := req1.URL.Query()
	q.Add("user_id", userID)
	req1.URL.RawQuery = q.Encode()

	c1 := &http.Client{}

	resp1, err := c1.Do(req1)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	body1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		panic(err)
	}

	token := string(body1)
	fmt.Println("Token: ", token)

	// Try to make a request to read the data
	fmt.Println(Green("Trying to fetch the data"))

	// Using HTTP bakery
	client := httpbakery.NewClient()

	// Decode the token and discharge anything
	fmt.Println(Blue("Discharging necessary caveats"))
	mac, err := macaroons.DecodeMacaroon(token)
	if err != nil {
		panic(err)
	}
	macs, err := client.DischargeAll(context.Background(), mac)
	if err != nil {
		panic(err)
	}

	fmt.Println(Green("Discharge succeeded, making actual request"))
	// Build and execute the actual request.
	//url := fmt.Sprintf("http://localhost:8080/api/acos/test/%s", acoID)
	url := fmt.Sprintf("http://localhost:3002/%s", acoID)

	httpbakery.SetCookie(client.Jar, mustParseURL(url), nil, macs)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(Green(string(body)))
}

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
