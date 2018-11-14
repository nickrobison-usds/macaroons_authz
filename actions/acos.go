package actions

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison/cms_authz/lib/auth/ca"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"github.com/nickrobison/cms_authz/lib/helpers"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
)

type idNamePair struct {
	ID   string
	Name string
}

var service *macaroons.Bakery

func init() {
	s, err := macaroons.NewBakery("http://localhost:8080/acos")
	if err != nil {
		log.Fatal(err)
	}
	service = s
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

func AcoShow(c buffalo.Context) error {
	acoID := c.Param("id")

	aco := models.ACO{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Eager().Find(&aco, acoID)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("aco", aco)

	// Add a binary helper
	c.Set("bytesToString", showBytes)

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
	condition := fmt.Sprintf("aco_id = %s", aco.ID)
	log.Debug(condition)
	mac, err := service.NewFirstPartyMacaroon([]string{condition})
	if err != nil {
		return err
	}
	b, err := macaroons.MacaroonToByteSlice(mac)
	if err != nil {
		return err
	}

	aco.Macaroon = b
	return nil
}

func AcoVerifyUser(c buffalo.Context) error {
	var requestData models.AcoUser
	err := c.Bind(&requestData)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something bad happened."))
	}

	log.Debugf("Verifying that user %s is a member of %s", requestData.UserID, requestData.ACOID)

	// Check that the association actually exists.
	tx := c.Value("tx").(*pop.Connection)

	var acoUser models.AcoUser

	err = tx.Where("aco_id = ?", requestData.ACOID).Where("user_id = ?", requestData.UserID).First(&acoUser)
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

// Returns a byte array as a string
func showBytes(s nulls.ByteSlice) string {
	return base64.URLEncoding.EncodeToString(s.ByteSlice)
}

/*
func dischargeUserCaveat(ctx context.Context, cav macaroon.Caveat, encodedCav []byte) (*macaroon.Macaroon, error) {

	log.Debug(cav.Id)
	log.Debug(cav.Location)

	mac, err := service.Discharge(macaroons.StrcmpChecker("user_id = test"), cav.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return mac, nil
}
*/

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

func DelegateACOToUser(acoID, userID uuid.UUID, tx *pop.Connection) error {
	// Create the intitial model
	link := models.AcoUser{
		ACOID:  acoID,
		UserID: userID,
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

	m, err := macaroons.MacaroonFromBytes(aco.Macaroon.ByteSlice)
	if err != nil {
		return err
	}

	// Add the caveats
	user_id := fmt.Sprintf("user_id = %s", userID.String())
	delegated, err := service.AddFirstPartyCaveats(m, []string{user_id})
	if err != nil {
		return err
	}

	mBinary, err := delegated.M().MarshalBinary()
	if err != nil {
		return err
	}

	link.Macaroon = mBinary

	return tx.Save(&link)
}
