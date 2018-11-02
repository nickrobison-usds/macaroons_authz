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
	C  string // Country
	ST string // State
	L  string // Locality
	O  string // OrganisationName
	OU string // OrganisationalUnitName
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

type CFSSLRequest struct {
	Token   []byte `json:"token,omitempty"`
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
		Hosts:        []string{"test.xample.com"},
		SerialNumber: "",
		Names: []Name{{
			C:  "US",
			L:  "Seattle",
			ST: "Washington",
			O:  name,
			OU: caName,
		}},
	}

	data, err := encodeCFSSLRequest(caRequest, "")
	if err != nil {
		return err
	}

	fmt.Printf("\n\nCert data: %s\n", data)

	// Send it over to CFSSL, key first, then Cert
	client := http.Client{}

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

	signedRequest := &CFSSLRequest{
		Request: signData,
	}

	token, err := signRequest("aaaaaaaaaaaaaaaa", signedRequest.Request)
	signedRequest.Token = token

	blob, err := json.Marshal(signedRequest)
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
		Contents: jsonCertResp.Result["private_key"].(string),
		IsBinary: false,
		Perms:    0664,
	})

	outs = append(outs, outputFile{
		FileName: caName + ".pem",
		Contents: jsonCertResp.Result["certificate"].(string),
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

// Sign a CFSSL request using the provided token
func signRequest(key string, data []byte) ([]byte, error) {
	keyHex, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	h := hmac.New(sha256.New, keyHex)
	h.Write(data)
	return h.Sum(nil), nil
}

// Encodes an API request in the format expected by CFSSL
func encodeCFSSLRequest(request interface{}, key string) ([]byte, error) {

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	cfsslRequest := &CFSSLRequest{
		Request: data,
	}

	// Encode it, if we need to
	if key != "" {
		token, err := signRequest(key, cfsslRequest.Request)
		if err != nil {
			return nil, err
		}
		cfsslRequest.Token = token
	}

	return json.Marshal(cfsslRequest)
}
