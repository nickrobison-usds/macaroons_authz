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

	"errors"

	"github.com/gobuffalo/envy"
	"github.com/mitchellh/mapstructure"
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
	Token   []byte      `json:"token,omitempty"`
	Request interface{} `json:"request"`
}

// responseMessage from CFSSL reporting errors or messages
type responseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CFSSLCertificateResponse is the result from CFSSL generating a new certifcate/private_key
type CFSSLCertificateResponse struct {
	PrivateKey  string `json:"private_key"`
	Certificate string `json:"certificate"`
	CSR         string `json:"certificate_request"`
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

var (
	cfsslURL string
	client   *http.Client
)

func init() {
	url, err := envy.MustGet("CFSSL_URL")
	if err != nil {
		panic(err)
	}

	cfsslURL = url
	client = &http.Client{}
}

// CreateCA - create a new CertificateAuthority for the given entity
func CreateCA(name string, caType Type) (CFSSLCertificateResponse, error) {
	// First, tell CFSSL to init a new CA.
	caName := fmt.Sprintf("%s_%s_ca", caType, name)

	caRequest := &CertificateRequest{
		CN: caName,
		KeyRequest: KeyRequest{
			Algo: "ecdsa",
			Size: 256,
		},
		SerialNumber: "",
		Names: []Name{{
			C:  "US",
			L:  "Seattle",
			ST: "Washington",
			O:  name,
			OU: caName,
		}},
	}

	fmt.Println("Doing the cert things")
	fmt.Println(caRequest)

	var jsonCertResp CFSSLCertificateResponse
	err := handleCFSSLOp(http.MethodPost, "newcert", caRequest, &jsonCertResp)
	if err != nil {
		return jsonCertResp, err
	}

	fmt.Println(jsonCertResp)

	// Now, pull out the data
	outs := []outputFile{}

	outs = append(outs, outputFile{
		FileName: caName + "-key.pem",
		Contents: jsonCertResp.PrivateKey,
		IsBinary: false,
		Perms:    0664,
	})

	outs = append(outs, outputFile{
		FileName: caName + ".pem",
		Contents: jsonCertResp.Certificate,
		IsBinary: false,
		Perms:    0664,
	})

	outs = append(outs, outputFile{
		FileName: caName + ".csr",
		Contents: jsonCertResp.CSR,
		Perms:    0664,
	})

	for _, e := range outs {
		err := ioutil.WriteFile(e.FileName, []byte(e.Contents), e.Perms)
		if err != nil {
			return jsonCertResp, err
		}
	}
	return jsonCertResp, nil
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

	//data, err := json.Marshal(request)
	//if err != nil {
	//	return nil, err
	//}

	cfsslRequest := &CFSSLRequest{
		Request: request,
	}

	// Encode it, if we need to
	//if key != "" {
	//	token, err := signRequest(key, cfsslRequest.Request)
	//	if err != nil {
	//		return nil, err
	//	}
	//	cfsslRequest.Token = token
	//}

	return json.Marshal(cfsslRequest)
}

func decodeCFSSLResponse(resp []byte, v interface{}) error {
	return json.Unmarshal(resp, v)
}

func handleCFSSLOp(method, resource string, data interface{}, output interface{}) error {
	encodedRequest, err := encodeCFSSLRequest(data, "")
	if err != nil {
		return err
	}

	uri := cfsslURL + "/api/v1/cfssl/" + resource

	// Execute HTTP operation
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(encodedRequest))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("\n\n%v\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	// Decode the response and check for any errors
	var response cfsslResponse
	err = decodeCFSSLResponse(body, &response)
	if err != nil {
		return err
	}

	fmt.Printf("Debug cfssl response: %v\n\n", response)

	if len(response.Errors) > 0 {
		err := response.Errors[0]
		return errors.New(err.Message)
	}

	// Now, decode the actual response
	err = mapstructure.Decode(response.Result, output)
	if err != nil {
		return err
	}
	return nil
}
