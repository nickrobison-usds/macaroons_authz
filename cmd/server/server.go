package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
	"gopkg.in/macaroon-bakery.v2/httpbakery"
	"io/ioutil"
	"net/http"
)

var b *bakery.Bakery

type keyResponse struct {
	K string `json:"k"`
}

func main() {
	fmt.Println("Starting test server")

	// Creating bakery
	loc := bakery.NewThirdPartyStore()
	// Get the key
	key, err := getKey("http://localhost:8080/api/vendors/.well-known/jwks.json")
	if err != nil {
		panic(err)
	}
	loc.AddInfo("http://localhost:8080/api/vendors/verify", bakery.ThirdPartyInfo{
		PublicKey: *key,
	})
	b = newBakery("http://localhost:3002", loc, nil)

	r := mux.NewRouter()
	r.HandleFunc("/{acoID}/token", TokenHandler)
	r.HandleFunc("/{acoID}", ResponseHandler)

	err = http.ListenAndServe(":3002", r)
	if err != nil {
		panic(err)
	}
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Get the user ID
	user_id := r.URL.Query().Get("user_id")
	aco_id := r.URL.Query().Get("acoID")
	fmt.Println("ACO ID: ", aco_id)

	// Add the third party caveat
	cav := checkers.Caveat{
		Location:  "http://localhost:8080/api/vendors/verify",
		Condition: fmt.Sprintf("user_id= %s", user_id),
	}

	mac, err := b.Oven.NewMacaroon(ctx, bakery.Version2, []checkers.Caveat{cav}, bakery.Op{"firstparty", "x"})
	if err != nil {
		panic(err)
	}

	json, err := mac.MarshalJSON()
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func ResponseHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	vars := mux.Vars(r)
	acoID := vars["acoID"]
	fmt.Println("Getting for ACO: ", acoID)

	// Get the Macarons from the request
	mSlice := httpbakery.RequestMacaroons(r)

	// Verify the macaroon
	ops, conds, err := b.Oven.VerifyMacaroon(ctx, mSlice[0])
	fmt.Println("Ops: ", ops)
	fmt.Println("Conds: ", conds)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Nope, can't do that."))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Successfully gathering data for ACO: %s", acoID)))
}

func newBakery(location string, locator bakery.ThirdPartyLocator, checker bakery.FirstPartyCaveatChecker) *bakery.Bakery {
	key, err := bakery.GenerateKey()
	if err != nil {
		panic(err)
	}
	return bakery.New(bakery.BakeryParams{
		Location: location,
		Locator:  locator,
		Key:      key,
		Checker:  nil,
	})
}

func getKey(location string) (*bakery.PublicKey, error) {

	publicKey := &bakery.PublicKey{}

	client := http.Client{}

	resp, err := client.Get(location)
	if err != nil {
		return publicKey, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return publicKey, err
	}

	var keyResp keyResponse

	err = json.Unmarshal(body, &keyResp)
	if err != nil {
		return publicKey, err
	}

	// Base 64 decode the key
	fmt.Println("Key: ", keyResp.K)
	// Add missing padding
	keyString, err := base64.URLEncoding.DecodeString(keyResp.K + "=")
	if err != nil {
		return nil, err
	}
	// Copy it into a fixed size array
	var key bakery.Key
	copy(key[:], keyString)
	publicKey.Key = key
	return publicKey, nil
}
