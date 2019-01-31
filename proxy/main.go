package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/envy"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/macaroons"
	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
)

var state map[string]string
var ps *macaroons.Bakery
var loginHost string

type userResponse struct {
	Email string `json:"email"`
	Sub   string `json:"sub"`
}

// dischargeResponse contains the response from a /discharge POST request.
type dischargeResponse struct {
	Macaroon *bakery.Macaroon `json:",omitempty"`
}

func main() {

	host := envy.Get("HOST", "http://localhost:5000")
	keyFile := envy.Get("KEY_FILE", "../user_keys.json")
	loginHost = envy.Get("LOGIN_HOST", "http://localhost:3000")

	// Setup our bakery

	// Read in the demo file
	f, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}

	key := &bakery.KeyPair{}
	err = json.Unmarshal(f, key)
	if err != nil {
		panic(err)
	}

	b, err := macaroons.NewBakery(host, createProxyChecker(), key)
	if err != nil {
		panic(err)
	}

	ps = b

	state = make(map[string]string)

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/discharge", dischargeHandler)

	fmt.Println("Listening")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		panic(err)
	}
}

func createProxyChecker() *checkers.Checker {
	c := checkers.New(nil)
	c.Namespace().Register("std", "")
	c.Register("email=", "std", macaroons.ContextCheck{"email"}.Check)

	return c
}

func dischargeHandler(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	fmt.Println("Token: ", token)
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Must provide OAuth token for discharge")
		return
	}
	mac := req.URL.Query().Get("id64")

	// Request the user info from the API
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/openid_connect/userinfo", loginHost), nil)
	if err != nil {
		handleError(w, err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Println(string(body))
	var user userResponse

	err = json.Unmarshal(body, &user)
	if err != nil {
		handleError(w, err)
		return
	}

	ctx := context.WithValue(req.Context(), "email", user.Email)

	// Discharge the Macaroon
	m2, err := ps.DischargeCaveatByID(ctx, mac, userEmailCaveatChecker())
	if err != nil {
		handleError(w, err)
		return
	}

	jsonResp, err := json.Marshal(&dischargeResponse{m2})
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Println("Writing macaroon")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		panic(err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err.Error())
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Proxy is running")
}

func userEmailCaveatChecker() bakery.ThirdPartyCaveatCheckerFunc {

	var caveats []checkers.Caveat
	return func(ctx context.Context, cav *bakery.ThirdPartyCaveatInfo) ([]checkers.Caveat, error) {
		email := ctx.Value("email")
		_, arg, err := checkers.ParseCaveat(string(cav.Condition))
		if err != nil {
			return caveats, err
		}

		fmt.Println("Checking that %s = %s", arg, email)
		if arg != email {
			return caveats, bakery.ErrPermissionDenied
		}

		return caveats, nil
	}
}
