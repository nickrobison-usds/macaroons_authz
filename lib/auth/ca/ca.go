package ca

import (
	"bytes"
	"fmt"
	"os"

	"encoding/hex"
	"encoding/json"

	"net/http"

	"io/ioutil"

	"crypto/hmac"
	"crypto/sha256"

	"github.com/gobuffalo/envy"
)

// The Type of the CA, either Vendor, ACO, or CMS
type Type string

// The following structs are taken from: https://github.com/cloudflare/cfssl/blob/master/csr/csr.go

// A Name contains the SubjectInfo fields.
type Name struct {
	C            string // Country
	ST           string // State
	L            string // Locality
	O            string // OrganisationName
	OU           string // OrganisationalUnitName
	SerialNumber string
}

// KeyRequest - Type and size of key to generate
type KeyRequest struct {
	Algo string `json:"algo"`
	Size int    `json:"size"`
}

// CertificateRequest
type CertificateRequest struct {
	CN           string
	Names        []Name     `json:"names" yaml:"names"`
	Hosts        []string   `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	KeyRequest   KeyRequest `json:"key,omitempty" yaml:"key,omitempty"`
	SerialNumber string     `json:"serialnumber,omitempty" yaml:"serialnumber,omitempty"`
}

type caConstraint struct {
	IsCA              bool `json:"is_ca"`
	MaxPathLength     int  `json:"max_path_len"`
	MaxPathLengthZero bool `json:"max_path_len_zero"`
}

type defaultSigningRequest struct {
	Usages       []string     `json:"usages"`
	Expiry       string       `json:"expiry"`
	CAConstraint caConstraint `json:"ca_constraint"`
	CRLURL       string       `json:"crl_url"`
}

type CsrRequest struct {
	CertificateRequest string `json:"certificate_request"`
}

type AuthSignRequest struct {
	Token   []byte `json:"token"`
	Request []byte `json:"request"`
}

// responseMessage from CFSSL reporting errors or messages
type responseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// cfsslResponse from CFSSL remote server
type cfsslResponse struct {
	Success  bool                   `json:"success"`
	Result   map[string]interface{} `json:"result"`
	Errors   []responseMessage      `json:"errors"`
	Messages []responseMessage      `json:"messages"`
}

// outputFile information
type outputFile struct {
	FileName string
	Contents string
	IsBinary bool
	Perms    os.FileMode
}

var cfsslURL string

func init() {
	url, err := envy.MustGet("CFSSL_URL")
	if err != nil {
		panic(err)
	}

	cfsslURL = url
}

// CreateCA - create a new CertificateAuthority for the given entity
func CreateCA(name string, caType Type) error {
	// First, tell CFSSL to init a new CA.
	caName := fmt.Sprintf("%s_%s_ca", caType, name)

	caRequest := &CertificateRequest{
		CN: caName,
		KeyRequest: KeyRequest{
			Algo: "ecdsa",
			Size: 256,
		},
		Hosts:        []string{},
		SerialNumber: "",
		Names: []Name{{
			C:  "US",
			L:  "Washington",
			ST: "District of Columbia",
			O:  name,
			OU: caName,
		}},
	}

	// Send it over to CFSSL, key first, then Cert
	data, err := json.Marshal(caRequest)
	if err != nil {
		return err
	}

	keyReq, err := http.NewRequest("POST", cfsslURL+"/api/v1/cfssl/newkey", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	fmt.Println("Doing the Key thing")

	client := &http.Client{}
	resp, err := client.Do(keyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Just print it out right now
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jsonKeyResp cfsslResponse

	err = json.Unmarshal(body, &jsonKeyResp)
	if err != nil {
		return err
	}

	fmt.Println(jsonKeyResp)

	// Get the Cert

	certReq, err := http.NewRequest("POST", cfsslURL+"/api/v1/cfssl/newcert", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	fmt.Println("Doing the cert things")
	fmt.Println(caRequest)

	certResp, err := client.Do(certReq)
	if err != nil {
		return err
	}
	defer certResp.Body.Close()

	certRespBody, err := ioutil.ReadAll(certResp.Body)
	if err != nil {
		return err
	}

	var jsonCertResp cfsslResponse
	err = json.Unmarshal(certRespBody, &jsonCertResp)
	if err != nil {
		return err
	}

	fmt.Println(jsonCertResp)

	// Now we have to sign the new intermediate
	fmt.Println("Signing Now")

	signData, err := json.Marshal(CsrRequest{
		CertificateRequest: jsonCertResp.Result["certificate_request"].(string)})
	if err != nil {
		return err
	}

	var debug map[string]interface{}
	err = json.Unmarshal(signData, &debug)
	if err != nil {
		return err
	}

	fmt.Println(debug)

	signRequest := &AuthSignRequest{
		Request: signData,
	}

	keyHex, err := hex.DecodeString("aaaaaaaaaaaaaaaa")
	if err != nil {
		return err
	}

	token, err := genToken(keyHex, signRequest.Request)
	signRequest.Token = token

	blob, err := json.Marshal(signRequest)
	if err != nil {
		return err
	}

	fmt.Println("Setting up req")

	signReq, err := http.NewRequest("POST", cfsslURL+"/api/v1/cfssl/authsign", bytes.NewReader(blob))
	if err != nil {
		return err
	}

	fmt.Println("Sending the sign request")
	signResp, err := client.Do(signReq)
	if err != nil {
		return err
	}
	defer signResp.Body.Close()

	respBody, err := ioutil.ReadAll(signResp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))

	var signRespData cfsslResponse

	err = json.Unmarshal(respBody, &signRespData)
	if err != nil {
		return err
	}

	fmt.Println(signRespData)

	// Now, pull out the data
	outs := []outputFile{}

	outs = append(outs, outputFile{
		FileName: caName + "-key.pem",
		Contents: jsonKeyResp.Result["private_key"].(string),
		IsBinary: false,
		Perms:    0664,
	})

	outs = append(outs, outputFile{
		FileName: caName + ".pem",
		Contents: signRespData.Result["certificate"].(string),
		IsBinary: false,
		Perms:    0664,
	})

	outs = append(outs, outputFile{
		FileName: caName + ".csr",
		Contents: jsonCertResp.Result["certificate_request"].(string),
		Perms:    0664,
	})

	for _, e := range outs {
		err := ioutil.WriteFile(e.FileName, []byte(e.Contents), e.Perms)
		if err != nil {
			return err
		}
	}

	return nil
}

func genToken(key, data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil), nil
}
