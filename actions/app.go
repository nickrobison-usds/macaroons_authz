package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr"
	"github.com/nickrobison/cms_authz/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator
var log logger.FieldLogger

func init() {
	log = logger.NewLogger("ACTIONS")
}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_cms_authz_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		// app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())
		app.Use(SetCurrentUser)

		app.GET("/", HomeHandler)

		app.GET("/routes", RouteHandler)

		// Auth handlers
		auth := app.Group("/auth")
		auth.POST("/login", ManualLogin)
		auth.GET("/logout", AuthDestroy)
		auth.GET("/{provider}", AuthLogin)
		auth.GET("/{provider}/callback", AuthCallback)

		api := app.Group("/api")
		api.Use(Authorize)
		api.Middleware.Skip(Authorize,
			AcosFind,
			AcoVerifyUser,
			AcoTest,
			AcoJWKS,
			UsersFind,
			UsersVerify,
			UsersTokenGet,
			VendorsVerify)

		// ACO Endpoints
		api.GET("/acos/create", RenderCreatePage)
		api.POST("/acos/create", AcosCreateACO)
		api.GET("/acos/delete/{id}", AcosDelete)
		api.GET("/acos/find", AcosFind)
		api.GET("/acos/index", AcosIndex)
		api.GET("/acos/list", AcosHeadIndex)
		api.GET("/acos/show/{id}", AcoShow)
		api.POST("/acos/verify", AcoVerifyUser)
		api.GET("/acos/test/{id}", AcoTest)
		api.GET("/acos/.well-known/jwks.json", AcoJWKS)

		// User Endpoints
		api.GET("/users/find", UsersFind)
		api.GET("/users/index", UsersIndex)
		api.GET("/users/show/{id}", UsersShow)
		api.POST("/users/create", UsersCreate)
		api.GET("/users/delete/{id}", UsersDelete)
		api.POST("/users/assign", UsersAssign)
		api.POST("/users/verify/discharge", UsersVerify)
		api.GET("/users/token/{user_id}/{entity_type}/{entity_id}", UsersTokenGet)

		// Vendor endpoints
		api.GET("/vendors/show/{id}", VendorsShow)
		api.GET("/vendors/index", VendorsIndex)
		api.GET("/vendors/list", VendorsList)
		api.GET("/vendors/create", VendorsCreate)
		api.POST("/vendors/assign", VendorsAssign)
		api.GET("/vendors/test/{id}", VendorsTest)
		api.POST("/vendors/{vendorID}/verify/discharge", VendorsVerify)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "test",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
