package actions

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/markbates/going/defaults"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
)

var usernames map[string]string

func init() {
	usernames = map[string]string{
		"nick@nick.com":  "password",
		"other@test.com": "test",
	}
	gothic.Store = App().SessionStore

	providerURL, err := envy.MustGet("PROVIDER_URL")
	if err != nil {
		panic(err)
	}

	discoveryURL := providerURL + "/.well-known/openid-configuration"

	oidp, err := openidConnect.New(envy.Get("CLIENT_ID", ""), os.Getenv("OPENIDCONNECT_SECRET"), fmt.Sprintf("%s%s", "http://localhost:8080", "/auth/login-gov/callback"), discoveryURL)
	if err != nil {
		panic(err)
	}

	goth.UseProviders(oidp)
}

// Custom login handler, because gothic.BeginAuthHandler doesn't work correctly.
func AuthLogin(c buffalo.Context) error {
	fmt.Println("Logging in")

	// Get the provider
	prov := c.Param("provider")
	provider, err := goth.GetProvider(prov)
	if err != nil {
		return errors.WithStack(err)
	}

	// Start the login
	state := generateNonce()

	sesh, err := provider.BeginAuth(state)
	if err != nil {
		errors.WithStack(err)
	}

	authURL, err := loginGovURL(sesh, state, "1")
	if err != nil {
		errors.WithStack(err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func AuthCallback(c buffalo.Context) error {
	gu, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("provider = ? and provider_id = ?", gu.Provider, gu.UserID)
	exists, err := q.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}
	u := &models.User{}
	if exists {
		if err = q.First(u); err != nil {
			return errors.WithStack(err)
		}
	}
	u.Name = defaults.String(gu.Name, gu.NickName)
	u.Provider = gu.Provider
	u.ProviderID = gu.UserID
	u.Email = nulls.NewString(gu.Email)
	if err = tx.Save(u); err != nil {
		return errors.WithStack(err)
	}

	c.Session().Set("current_user_id", u.ID)
	if err = c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "You have been logged in")
	return c.Redirect(302, "/")
}

func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out")
	return c.Redirect(302, "/")
}

func generateNonce() string {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic(err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// Generate the login URL
func loginGovURL(session goth.Session, state string, loaNum string) (string, error) {
	urlStr, err := session.GetAuthURL()
	if err != nil {
		return "", err
	}

	authURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	params := authURL.Query()
	params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/"+loaNum)
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

func ManualLogin(c buffalo.Context) error {
	r := c.Request()
	fmt.Println(r.Form)

	email := r.FormValue("email")

	password, ok := usernames[email]
	if !ok {
		return errors.WithStack(errors.New("Incorrect user"))
	}

	if r.FormValue("password") != password {
		return errors.WithStack(errors.New("Bad password"))
	}

	fmt.Println("Logged")
	c.Session().Set("user_id", email)

	return c.Redirect(302, "/")

}

func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			if err := tx.Find(u, uid); err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
