package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

// Stores token information.
// See: https://developers.login.gov/oidc/#token-response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

// Fetch self-signed JWT token which corresponds to OAuth session
func FetchToken(c buffalo.Context) (tokenResponse, error) {
	var tr tokenResponse

	providerURL := envy.Get("PROVIDER_URL", "")
	tokenURL := providerURL + "/api/openid_connect/token"
	// Get the clientId from the context
	clientId := c.Value("client_id").(string)

	clientAssertion, err := generateJWT(tokenURL, clientId)
	if err != nil {
		return tr, err
	}

	queryParams := c.Request().URL.Query()
	code := queryParams["code"][0]

	tokenParams := url.Values{}
	tokenParams.Set("client_assertion", clientAssertion)
	tokenParams.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	tokenParams.Set("code", code)
	tokenParams.Set("grant_type", "authorization_code")

	// Request the token
	resp, err := http.PostForm(tokenURL, tokenParams)
	if err != nil {
		return tr, err
	}

	// Parse response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tr, err
	}

	err = json.Unmarshal(body, &tr)
	if err != nil {
		return tr, err
	}

	return tr, nil
}

func FetchUserInfo(tr tokenResponse) (goth.User, error) {

	provider, err := goth.GetProvider("login-gov")
	if err != nil {
		return goth.User{}, err
	}

	session := openidConnect.Session{
		AccessToken: tr.AccessToken,
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(tr.ExpiresIn)),
		IDToken:     tr.IDToken,
	}

	return provider.FetchUser(&session)
}

// GenerateNonce creates a random value for use in Login.gov authentication
func GenerateNonce() string {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic(err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func GetLoginGovUser(c buffalo.Context) (goth.User, error) {

	tokenResponse, err := FetchToken(c)
	if err != nil {
		return goth.User{}, err
	}
	return FetchUserInfo(tokenResponse)
}

// Private methods

func generateJWT(tokenURL, clientId string) (string, error) {
	pem, err := ioutil.ReadFile("login-gov/macdemo.key")
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return "", err
	}

	// Generate a new token

	const sessionExpiryInMinutes = 120

	claims := &jwt.StandardClaims{
		Issuer:    clientId,
		Subject:   clientId,
		Audience:  tokenURL,
		Id:        GenerateNonce(),
		ExpiresAt: time.Now().Add(time.Minute * sessionExpiryInMinutes).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	jwt, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
