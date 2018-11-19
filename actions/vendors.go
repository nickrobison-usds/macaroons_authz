package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison/cms_authz/lib/auth/ca"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"github.com/nickrobison/cms_authz/lib/helpers"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
)

const vendorURI = "http://localhost:8080/api/vendors"

type vendorAssignRequest struct {
	VendorID string `form:"vendorID"`
	EntityId string `form:"entityOptions"`
}

var vs *macaroons.Bakery

func init() {
	b, err := macaroons.NewBakery(vendorURI, createVendorChecker(), models.DB, nil)
	if err != nil {
		log.Fatal(err)
	}

	vs = b
}

func createVendorChecker() *checkers.Checker {
	c := checkers.New(nil)
	c.Namespace().Register("std", "")
	c.Register("vendor_id=", "std", macaroons.ContextCheck{"vendor_id"}.StrCheck)

	return c
}

// VendorsAssign assigns a vendor to a given ACO
func VendorsAssign(c buffalo.Context) error {

	req := vendorAssignRequest{}

	err := c.Bind(&req)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Debugf("Assigning vendor %s to ACO %s.\n", req.VendorID, req.EntityId)

	tx := c.Value("tx").(*pop.Connection)

	err = AssignVendorToACO(helpers.UUIDOfString(req.EntityId), helpers.UUIDOfString(req.VendorID), tx)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/api/vendors/index")
}

// VendorsCreate default implementation.
func VendorsCreate(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/vendors/create.html"))
}

func CreateVendorCertificates(vendor *models.Vendor) error {

	// Generate CA
	cert, err := ca.CreateCA(vendor.Name, "vendor")
	if err != nil {
		return err
	}

	// Not used yet
	_, err = TransformCFSSLResponse(vendor.ID, &cert)
	if err != nil {
		return err
	}

	// Create a Macaroon for the Vendor
	condition := fmt.Sprintf("vendor_id= %s", vendor.ID)
	log.Debug(condition)

	mac, err := vs.NewFirstPartyMacaroon([]string{condition})
	if err != nil {
		return err
	}

	b, err := macaroons.MacaroonToByteSlice(mac)
	if err != nil {
		return err
	}

	vendor.Macaroon = b.ByteSlice
	return nil
}

// AssignUserToVendor creates a delegated macaroon from the vendor, to the appropriate user.
func AssignUserToVendor(vendorID, userID uuid.UUID, tx *pop.Connection) error {
	// Create the initial model
	link := models.VendorUser{
		VendorID: vendorID,
		UserID:   userID,
	}

	// Get the macaroon from the vendor
	vendor := models.Vendor{}

	err := tx.Select("macaroon").Where("id = ?",
		vendorID.String()).First(&vendor)
	if err != nil {
		return err
	}

	// Generate the delegating Macaroon
	m, err := macaroons.MacaroonFromBytes(vendor.Macaroon)
	if err != nil {
		return err
	}

	// Add the caveats
	userId := fmt.Sprintf("user_ID= %s", userID.String())
	delegated, err := vs.AddFirstPartyCaveats(m, []string{userId})
	if err != nil {
		return err
	}

	dBinary, err := delegated.M().MarshalBinary()
	if err != nil {
		return err
	}

	link.Macaroon = dBinary

	return tx.Save(&link)
}

// VendorsIndex default implementation.
func VendorsIndex(c buffalo.Context) error {

	vendors := []models.Vendor{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.All(&vendors); err != nil {
		return errors.WithStack(err)
	}

	c.Set("vendors", vendors)
	return c.Render(200, r.HTML("api/vendors/index.html"))
}

func VendorsList(c buffalo.Context) error {

	vendors := []models.Vendor{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.All(&vendors); err != nil {
		return errors.WithStack(err)
	}

	values := []idNamePair{}

	for _, vendor := range vendors {
		values = append(values, idNamePair{Name: vendor.Name, ID: vendor.StringID()})
	}

	return c.Render(200, r.JSON(&values))
}

// VendorsShow default implementation.
func VendorsShow(c buffalo.Context) error {

	vendorID := c.Param("id")

	vendor := &models.Vendor{
		ID: helpers.UUIDOfString(vendorID),
	}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Eager().Find(vendor, vendor.ID); err != nil {
		return errors.WithStack(err)
	}

	c.Set("vendor", vendor)

	return c.Render(200, r.HTML("api/vendors/show.html"))
}

// VendorsTest tests whether a delegated vendor token can be used to access resources.
func VendorsTest(c buffalo.Context) error {
	log.Debug("Attempting to verify test vendor macaroon")

	vendorID := c.Param("id")
	token := c.Param("token")

	m, err := macaroons.DecodeMacaroon(token)
	if err != nil {
		return errors.WithStack(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "vendor_id", vendorID)
	ctx = context.WithValue(ctx, "user_id", "test-user")

	err = vs.VerifyMacaroon(ctx, m)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.String("success! %s", vendorID))
}
