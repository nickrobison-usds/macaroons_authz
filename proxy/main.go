package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var state map[string]string

func main() {

	state = make(map[string]string)

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/discharge", dischargeHandler)

	fmt.Println("Listening")
	err := http.ListenAndServe(":5000", r)
	if err != nil {
		panic(err)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Proxy is running")
}

func dischargeHandler(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	fmt.Println("Token: ", token)
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Must provide OAuth token for discharge")
		return
	}

	// Request the user info from the API
	req, err := http.NewRequest("GET", "http://localhost:3000/api/openid_connect/userinfo", nil)
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Discharged")

}

func handleError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err.Error())
	return
}
