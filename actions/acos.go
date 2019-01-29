package actions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/ca"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/macaroons"
	"github.com/nickrobison-usds/macaroons_authz/lib/helpers"
	"github.com/nickrobison-usds/macaroons_authz/models"
	"github.com/pkg/errors"
	"github.com/rakutentech/jwk-go/jwk"
	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
	"gopkg.in/macaroon-bakery.v2/httpbakery"
)

var acoURI = "http://localhost:8080/api/acos/verify"

type idNamePair struct {
	ID   string
	Name string
}

var as *macaroons.Bakery

func init() {
	// Read in the demo file
	f, err := ioutil.ReadFile("./user_keys.json")
	if err != nil {
		log.Fatal(err)
	}

	key := &bakery.KeyPair{}
	err = json.Unmarshal(f, key)
	if err != nil {
		log.Fatal(err)
	}
	s, err := macaroons.NewBakery(acoURI, createACOCheckers(), models.DB, key)
	if err != nil {
		log.Fatal(err)
	}
	as = s
}

func AcosCreateACO(c buffalo.Context) error {

	fmt.Println(c.Request())

	aco := models.ACO{}

	aco.ID = helpers.MustGenerateID()

	err := c.Bind(&aco)
	if err != nil {
		return errors.WithStack(err)
	}
	err = CreateACOCertificates(&aco)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("\n\n\nACO: %v\n\n\n", aco)

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Create(aco); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/api/acos/index")
}

func AcosDelete(c buffalo.Context) error {

	acoID := c.Param("id")
	uuidString, err := uuid.FromString(acoID)
	if err != nil {
		return errors.WithStack(err)
	}

	user := models.ACO{
		ID: uuidString,
	}
	tx := c.Value("tx").(*pop.Connection)
	err = tx.Destroy(&user)
	if err != nil {
		c.Flash().Add("danger", fmt.Sprintf("Cannot delete ACO: %s", acoID))
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Deleted")
	return c.Redirect(302, "/api/acos/index")

}

func AcoDischargeMacaroon(c buffalo.Context) error {
	// Retrieve the token from the request
	token := c.Param("id64")

	// Get the ACO ID from the path
	acoID := c.Param("id")

	log.Debug("Before ctx: ", c.Request().Context())
	log.Debug("user param:", c.Param("user_id"))
	log.Debug("Req: ", c.Request())
	// Add the user ID to the context
	ctx := context.WithValue(c.Request().Context(), "user_id", c.Param("user_id"))
	log.Debug("after ctx: ", ctx)

	// Vendor?
	vendorID := c.Param("vendor_id")
	log.Debug("Vendor ID:", vendorID)

	// If it comes in as id, we need to translate it to base64 encoding
	if token == "" {
		token = base64.URLEncoding.EncodeToString([]byte(c.Param("id")))
	}

	log.Debug("Token: ", token)

	mac, err := us.DischargeCaveatByID(ctx, token, userAssociatedChecker(c.Value("tx").(*pop.Connection), helpers.UUIDOfString(acoID)))
	if err != nil {
		log.Debug(err)
		// Do a string compare on the permission denied type
		// This is gross, but I'm not sure how to get the underlying causer
		if err.Error() == "permission denied" {
			errMsg := fmt.Sprintf("Not authorized to retrieve data for ACO %s", acoID)
			log.Error("Unauthorized: ", err)
			return c.Render(http.StatusUnauthorized, r.String(errMsg))
		}
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(&dischargeResponse{mac}))
}

// AcosFind looks up an ACO ID via the given parameter
func AcosFind(c buffalo.Context) error {
	nameString := c.Param("name")
	if nameString == "" {
		return c.Render(http.StatusBadRequest, r.String("Cannot have a blank query name."))
	}

	aco := models.ACO{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Where("name = ?", nameString).First(&aco); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, r.String(aco.StringID()))
}

// acoIndex default implementation.
func AcosIndex(c buffalo.Context) error {

	acos := []models.ACO{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.All(&acos); err != nil {
		return errors.WithStack(err)
	}

	c.Set("acos", acos)
	return c.Render(200, r.HTML("api/acos/index.html"))
}

func AcosHeadIndex(c buffalo.Context) error {
	acos := []models.ACO{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Select("id, name").All(&acos); err != nil {
		return errors.WithStack(err)
	}

	values := []idNamePair{}

	for _, aco := range acos {
		values = append(values, idNamePair{Name: aco.Name, ID: aco.StringID()})
	}

	return c.Render(200, r.JSON(&values))

}

// AcoJWKS returns the jwks.json for the given third-party
func AcoJWKS(c buffalo.Context) error {
	pub := as.GetPublicKey()
	spec := jwk.NewSpec(pub)

	spec.KeyID = "1"
	spec.Algorithm = "ES256"
	spec.Use = "enc"

	return c.Render(http.StatusOK, r.JSON(spec))
}

func AcoShow(c buffalo.Context) error {
	acoID := c.Param("id")

	aco := models.ACO{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Eager().Find(&aco, acoID)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("aco", aco)

	return c.Render(http.StatusOK, r.HTML("/api/acos/show.html"))
}

// CreateACOCertificates sets up the certs, keys and macaroons for the given ACO
func CreateACOCertificates(aco *models.ACO) error {

	// Generate CA
	cert, err := ca.CreateCA(aco.Name, "aco")
	if err != nil {
		return err
	}

	parsed, err := TransformCFSSLResponse(aco.ID, &cert)
	if err != nil {
		return err
	}

	aco.Certificate = parsed

	// Now do the macaroon

	// Create a Macaroon for the ACO
	condition := fmt.Sprintf("aco_id= %s", aco.ID)
	log.Debug(condition)
	mac, err := as.NewFirstPartyMacaroon([]string{condition})
	if err != nil {
		return err
	}
	b, err := mac.M().MarshalBinary()
	if err != nil {
		return err
	}

	aco.Macaroon = b
	return nil
}

func AcoTest(c buffalo.Context) error {
	log.Debug("Trying to test that it works.")
	acoId := c.Param("id")

	// Get the macaroons from the request
	// Right now, we only need the first set.
	mSlice := httpbakery.RequestMacaroons(c.Request())

	log.Debug("Verifying macaroon")

	// Verify
	// Gen context
	ctx := context.WithValue(c.Request().Context(), "aco_id", acoId)

	err := as.VerifyMacaroons(ctx, mSlice[0])
	if err != nil {
		log.Errorf("Auth error: %s", err.Error())
		return c.Render(http.StatusUnauthorized, r.String("Unauthorized"))
	}

	return c.Render(200, r.String("Successfully accessed data for: %s", acoId))
}

func createACOCheckers() *checkers.Checker {
	c := checkers.New(nil)
	c.Namespace().Register("std", "")
	c.Register("entity_id=", "std", macaroons.CMSAssociationCheck{
		ContextKey:        "aco_id",
		AssociationTable:  "aco_users",
		AssociationColumn: "entity_id",
		DB:                models.DB,
	}.Check)
	c.Register("aco_id=", "std", macaroons.ContextCheck{Key: "aco_id"}.Check)
	/*c.Register("user_id=", "std", macaroons.CMSAssociationCheck{
		ContextKey:        "aco_id",
		AssociationTable:  "aco_users",
		AssociationColumn: "user_id",
		DB:                models.DB,
	}.Check)
	*/

	return c
}

func AcoVerifyUser(c buffalo.Context) error {
	var requestData models.AcoUser
	err := c.Bind(&requestData)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something bad happened."))
	}

	log.Debugf("Verifying that user %s is a member of %s", requestData.EntityID, requestData.ACOID)

	// Check that the association actually exists.
	tx := c.Value("tx").(*pop.Connection)

	var acoUser models.AcoUser

	err = tx.Where("aco_id = ?", requestData.ACOID).Where("user_id = ?", requestData.EntityID).First(&acoUser)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something went wrong: %s", err.Error()))
	}

	// If it exists, discharge it
	_, err = macaroons.MacaroonFromBytes(requestData.Macaroon)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something went wrong: %s", err.Error()))
	}

	// Validate  it?

	return c.Render(http.StatusOK, r.String("ok"))
}

// TransformCFSSLResponse converts a ca.CFSSLCertificateResponse into a models.Certificate
func TransformCFSSLResponse(id uuid.UUID, cert *ca.CFSSLCertificateResponse) (models.Certificate, error) {
	fmt.Println(cert)

	parsed, err := ca.ParseCFSSLResponse(cert)
	if err != nil {
		return models.Certificate{}, errors.WithStack(err)
	}

	encCert, err := parsed.EncodeCertificate()
	if err != nil {
		return models.Certificate{}, errors.WithStack(err)
	}

	priv, err := parsed.EncodePrivateKey()
	if err != nil {
		return models.Certificate{}, errors.WithStack(err)
	}

	fmt.Printf("ACO ID: %s\n", id)

	acoCert := models.Certificate{
		ACOID:       id,
		Key:         priv,
		Certificate: encCert,
		SHA:         parsed.SHA,
	}

	return acoCert, nil
}

func AssignUserToACO(acoID, userID uuid.UUID, tx *pop.Connection) error {
	// Create the intitial model
	link := models.AcoUser{
		ACOID:    acoID,
		EntityID: userID,
		IsUser:   true,
	}

	// Get the Macaroon from the ACO
	aco := models.ACO{}

	err := tx.Select("macaroon").Where("id = ?",
		acoID.String()).First(&aco)
	if err != nil {
		return err
	}

	// Generate a delegating Macaroon
	// We need a third party to attest that the user is who they say they are.

	// Decode macaroon
	m, err := macaroons.MacaroonFromBytes(aco.Macaroon)
	if err != nil {
		return err
	}

	entityCaveat := fmt.Sprintf("entity_id= %s", userID.String())
	delegated, err := as.AddFirstPartyCaveats(m, []string{entityCaveat})
	if err != nil {
		return err
	}

	// Add a third party caveat
	userCaveat := fmt.Sprintf("user_id= %s", userID.String())
	d1, err := as.AddThirdPartyCaveat(delegated, "http://localhost:8080/api/users/verify", []string{userCaveat})
	if err != nil {
		return err
	}

	log.Debug("User assigned token:", d1)

	mBinary, err := d1.M().MarshalBinary()
	if err != nil {
		return err
	}

	link.Macaroon = mBinary

	return tx.Save(&link)
}

func AssignVendorToACO(acoID, vendorID uuid.UUID, tx *pop.Connection) error {
	// Create the initial model
	link := models.AcoUser{
		ACOID:    acoID,
		EntityID: vendorID,
		IsUser:   false,
	}

	// Get the ACO's Macaroon
	aco := models.ACO{}

	err := tx.Select("macaroon").Where("id = ?", acoID.String()).First(&aco)
	if err != nil {
		return err
	}

	// Generate a delegating Macaroon

	m, err := macaroons.MacaroonFromBytes(aco.Macaroon)
	if err != nil {
		return err
	}

	// Restrict macaroon to only that Vendor

	vendorCaveat := fmt.Sprintf("entity_id= %s", vendorID.String())

	// Verify that the vendor is known to the ACO
	d1, err := as.AddFirstPartyCaveats(m, []string{vendorCaveat})
	if err != nil {
		return err
	}

	mBinary, err := d1.M().MarshalBinary()
	if err != nil {
		return err
	}

	link.Macaroon = mBinary

	return tx.Save(&link)
}

func RenderCreatePage(c buffalo.Context) error {
	aco := models.ACO{}

	c.Set("aco", aco)
	return c.Render(http.StatusOK, r.HTML("api/acos/create.html"))
}

func userAssociatedChecker(db *pop.Connection, acoID uuid.UUID) bakery.ThirdPartyCaveatCheckerFunc {
	return func(ctx context.Context, cav *bakery.ThirdPartyCaveatInfo) ([]checkers.Caveat, error) {

		var caveats []checkers.Caveat
		_, entityID, err := checkers.ParseCaveat(string(cav.Condition))
		if err != nil {
			return nil, err
		}
		log.Debugf("Checking that %s is associated with ACO %s", entityID, acoID)

		var user models.AcoUser

		err = db.Where("entity_id = ? AND aco_id = ?", helpers.UUIDOfString(entityID), acoID).First(&user)
		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows {
				return caveats, bakery.ErrPermissionDenied
			}
			return caveats, err
		}

		// If we're a vendor, add a third party caveat requiring the vendor to verify that the user is valid
		proxyHost := envy.Get("LOGIN_PROXY", "")

		if user.IsUser && proxyHost != "" {
			log.Debug("Setting login host: ", proxyHost)

			// Fetch the user, so we can get its openIDConnec
			// Add a caveat requiring the proxy service to verify with oauth.
			cav, err := createUserCaveat(&user, proxyHost, db)
			if err != nil {
				return caveats, err
			}
			caveats = append(caveats, cav)
		} else {
			// Get the user from the context
			userID := ctx.Value("user_id").(string)
			if userID == "" {
				return caveats, fmt.Errorf("Need userID when discharging as a vendor.")
			}
			log.Debug("Entity is a vendor, adding additional caveat for user: ", userID)
			caveats = append(caveats,
				checkers.Caveat{
					Location:  fmt.Sprintf("http://localhost:8080/api/vendors/%s/verify", entityID),
					Condition: fmt.Sprintf("user_id= %s", userID),
				})
		}

		return caveats, nil
	}
}

func createUserCaveat(acoUser *models.AcoUser, proxy string, db *pop.Connection) (checkers.Caveat, error) {
	var caveat checkers.Caveat
	var user models.User

	err := db.Where("id = ?", acoUser.EntityID).First(&user)
	if err != nil {
		return caveat, err
	}

	caveat.Location = proxy
	caveat.Condition = fmt.Sprintf("user_id= %s", user.ProviderID)
	return caveat, nil
}
