package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/ca"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/macaroons"
	"github.com/nickrobison-usds/macaroons_authz/lib/helpers"
	"github.com/nickrobison-usds/macaroons_authz/models"
	"github.com/pkg/errors"
	"github.com/rakutentech/jwk-go/jwk"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
	"gopkg.in/macaroon-bakery.v2/httpbakery"
)

var acoURI = "http://localhost:8080/api/acos"

type idNamePair struct {
	ID   string
	Name string
}

var as *macaroons.Bakery

func init() {
	s, err := macaroons.NewBakery(acoURI, createACOCheckers(), models.DB, nil)
	if err != nil {
		log.Fatal(err)
	}
	as = s
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

func RenderCreatePage(c buffalo.Context) error {
	aco := models.ACO{}

	c.Set("aco", aco)
	return c.Render(http.StatusOK, r.HTML("api/acos/create.html"))
}

// AcoJWKS returns the jwks.json for the given third-party
func AcoJWKS(c buffalo.Context) error {
	priv := as.GetPrivateKey()
	spec := jwk.NewSpec(priv)

	spec.KeyID = "1"
	spec.Algorithm = "ES256"
	spec.Use = "enc"

	log.Debug("Spec: ", spec)

	json, err := spec.MarshalJSON()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Debug(json)
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
		return c.Render(http.StatusInternalServerError, r.String("Something wen't wrong: %s", err.Error()))
	}

	// If it exists, discharge it
	_, err = macaroons.MacaroonFromBytes(requestData.Macaroon)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something wen't wrong: %s", err.Error()))
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
