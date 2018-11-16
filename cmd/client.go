package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gobuffalo/envy"
	. "github.com/logrusorgru/aurora"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"gopkg.in/macaroon-bakery.v2/httpbakery"
	macaroon "gopkg.in/macaroon.v2"
)

var acoID = "ba62259e-8083-4a2e-8de4-bd55622a9dd4"

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

	// Can we try to use the httpbakery?

	mac, err := macaroons.DecodeMacaroon(token)
	if err != nil {
		panic(err)
	}

	client := httpbakery.NewClient()

	macs, err := client.DischargeAll(context.Background(), mac)
	if err != nil {
		panic(err)
	}

	mBinary, err := macs[0].MarshalBinary()
	if err != nil {
		panic(err)
	}

	// Build and execute the actual request.
	url := fmt.Sprintf("http://localhost:8080/api/acos/test/%s", acoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Macaroons", macaroons.EncodeMacaroon(mBinary))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	/*


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

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
