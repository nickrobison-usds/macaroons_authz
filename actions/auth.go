package actions

import (
	"fmt"
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
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth"
	"github.com/nickrobison-usds/macaroons_authz/models"
	"github.com/pkg/errors"
)

var (
	usernames map[string]string
)

func init() {

	// Get rid of this
	usernames = map[string]string{
		"nick@nick.com":  "password",
		"other@test.com": "test",
	}
	gothic.Store = App().SessionStore

	providerURL := envy.Get("PROVIDER_URL", "")

	discoveryURL := providerURL + "/.well-known/openid-configuration"

	providers := []goth.Provider{}

	oidp, err := openidConnect.New(envy.Get("CLIENT_ID", ""), os.Getenv("OPENIDCONNECT_SECRET"), fmt.Sprintf("%s%s", "http://localhost:8080", "/auth/login-gov/callback"), discoveryURL)
	if err == nil {
		oidp.SetName("login-gov")
		providers = append(providers, oidp)
	} else {
		log.Warn("Not enabling Loging.gov: ", err)
	}

	// Github provider
	gidb := github.New(envy.Get("GITHUB_KEY", ""), envy.Get("GITHUB_SECRET", ""), "http://localhost:8080/auth/github/callback")
	providers = append(providers, gidb)

	goth.UseProviders(providers...)
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

	log.Debug("Provider:", prov)

	// Special handling of Login.Gov, because it has some funkiness to it.
	if prov == "login-gov" {
		// Start the login
		state := auth.GenerateNonce()

		sesh, err := provider.BeginAuth(state)
		if err != nil {
			return errors.WithStack(err)
		}

		authURL, err := loginGovURL(sesh, state, "1")
		if err != nil {
			return errors.WithStack(err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}

	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil

}

func AuthCallback(c buffalo.Context) error {

	gu, err := getLoginUser(c)
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("email = ?", gu.Email)
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

// AuthDestroy logs out the user and clears the session
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out")
	return c.Redirect(302, "/")
}

func getLoginUser(c buffalo.Context) (goth.User, error) {
	var gu goth.User
	var err error
	// Special handling of Login.gov auth
	prov := c.Param("provider")

	if prov == "login-gov" {
		// Set the client ID in the context
		c.Set("client_id", envy.Get("CLIENT_ID", ""))
		gu, err = auth.GetLoginGovUser(c)
		if err != nil {
			return gu, err
		}
	} else {
		gu, err = gothic.CompleteUserAuth(c.Response(), c.Request())
		if err != nil {
			return gu, err
		}
	}
	return gu, nil
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
	if !ok || password != r.FormValue("password") {
		c.Flash().Add("auth", "Incorrect username or password")
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
