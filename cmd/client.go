package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gobuffalo/envy"
	. "github.com/logrusorgru/aurora"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"gopkg.in/macaroon-bakery.v2/httpbakery"
	macaroon "gopkg.in/macaroon.v2"
)

var acoID = "f92af150-7c69-4f56-958e-9f5fd529c57c"

func main() {
	token, err := envy.MustGet("TOKEN")
	if err != nil {
		panic(err)
	}
	fmt.Println(Green("Starting up"))

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
	url := fmt.Sprintf("http://localhost:8080/api/acos/test/%s", acoID)

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
