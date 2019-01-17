package actions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/ca"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/macaroons"
	"github.com/nickrobison-usds/macaroons_authz/lib/helpers"
	"github.com/nickrobison-usds/macaroons_authz/models"
	"github.com/pkg/errors"
	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
)

const vendorURI = "http://localhost:8080/api/vendors/verify"

type vendorAssignRequest struct {
	VendorID string `form:"vendorID"`
	EntityId string `form:"entityOptions"`
}

var vs *macaroons.Bakery

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

	b, err := macaroons.NewBakery(vendorURI, createVendorChecker(), models.DB, key)
	if err != nil {
		log.Fatal(err)
	}

	vs = b
}

func createVendorChecker() *checkers.Checker {
	c := checkers.New(nil)
	c.Namespace().Register("std", "")
	c.Register("entity_id=", "std", macaroons.ContextCheck{"entity_id"}.Check)

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

func VendorsFind(c buffalo.Context) error {
	nameString := c.Param("name")
	if nameString == "" {
		return c.Render(http.StatusBadRequest, r.String("Cannot have blank query name"))
	}

	vendor := models.Vendor{}

	tx := c.Value("tx").(*pop.Connection)
	err := tx.Select("id").Where("name = ?", nameString).First(&vendor)
	if err != nil {
		log.Error(err)
		if errors.Cause(err) == sql.ErrNoRows {
			return c.Render(http.StatusNotFound, r.String(fmt.Sprintf("Cannot find vendor with name %s", nameString)))
		}
		return c.Render(http.StatusInternalServerError, r.String("Something went wrong."))
	}

	return c.Render(http.StatusOK, r.String(vendor.ID.String()))
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
	log.Debug("Creating vendor macaroon.")
	condition := fmt.Sprintf("entity_id= %s", vendor.ID)
	log.Debug(condition)
	mac, err := vs.NewFirstPartyMacaroon([]string{condition})
	if err != nil {
		return err
	}
	// Add a third party caveat to verify that the vendor is assigned to the aco
	/*
		mac, err := vs.NewThirdPartyMacaroon(context.Background(), "http://localhost:8080/api/acos/verify", []string{condition})
		if err != nil {
			return err
		}
	*/

	log.Debug("Vendor thing:", mac)

	b, err := mac.M().MarshalBinary()
	if err != nil {
		return err
	}

	vendor.Macaroon = b
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
	userId := fmt.Sprintf("user_id= %s", userID.String())
	// Use a third party caveat, because we need to have other APIs verify the user list.

	verifyString := fmt.Sprintf("http://localhost:8080/api/vendors/%s/verify", vendorID.String())
	delegated, err := vs.AddThirdPartyCaveat(m, verifyString, []string{userId})
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

	vendorID := c.Param("vendorID")
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

// VendorsVerify verifies that a given user is a member of the Vendor
func VendorsVerify(c buffalo.Context) error {
	log.Debug("Test:", c.Request())
	token := c.Param("id64")

	log.Debug("Token: ", token)
	log.Debug("Checking user association for vendor: ", c.Param("vendorID"))
	// Set the vendor_id
	ctx := context.WithValue(c.Request().Context(), "vendor_id", c.Param("vendorID"))

	// Decode the caveat, and keep going
	mac, err := vs.DischargeCaveatByID(ctx, token, vendorUserIDCaveatChecker(c.Value("tx").(*pop.Connection)))
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(&dischargeResponse{mac}))
}

func vendorUserIDCaveatChecker(db *pop.Connection) bakery.ThirdPartyCaveatCheckerFunc {

	return func(ctx context.Context, cav *bakery.ThirdPartyCaveatInfo) ([]checkers.Caveat, error) {

		// Get the vendor ID from the context
		vendorIDString := ctx.Value("vendor_id").(string)
		vendorID := helpers.UUIDOfString(vendorIDString)

		var caveats []checkers.Caveat
		log.Debug("In the Vendor ID checker")
		log.Debug(string(cav.Condition))
		_, arg, err := checkers.ParseCaveat(string(cav.Condition))
		if err != nil {
			return caveats, err
		}

		var vendor models.VendorUser

		// Getting from the DB
		ID := helpers.UUIDOfString(arg)

		err = db.Where("vendor_id = ? and user_id = ?", vendorID, ID).First(&vendor)
		if err != nil {
			log.Error(err)
			if errors.Cause(err) == sql.ErrNoRows {
				return caveats, bakery.ErrPermissionDenied
			}
			return caveats, err
		}

		log.Debug("Found user:", vendor)

		// The user is known to us, but are they valid?
		// Need to verify

		log.Debug("Adding user login caveat")

		return []checkers.Caveat{checkers.Caveat{
			Location:  "http://localhost:8080/api/users/verify",
			Condition: fmt.Sprintf("user_id= %s", vendor.UserID.String()),
			Namespace: checkers.StdNamespace,
		}}, nil
	}

}
