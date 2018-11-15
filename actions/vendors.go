package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison/cms_authz/lib/auth/ca"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
)

const vendorURI = "http://localhost:8080/api/vendors"

var vs *macaroons.Bakery

func init() {
	b, err := macaroons.NewBakery(vendorURI, createVendorChecker())
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

// DelegateUserToVendor creates a delegated macaroon from the vendor, to the appropriate user.
func DelegateUserToVendor(vendorID, userID uuid.UUID, tx *pop.Connection) error {
	// Create the initial model
	link := models.VendorUser{
		VendorID: vendorID,
		UserID:   userID,
	}

	// Get the macaroon from the vendor
	vendor := models.ACO{}

	err := tx.Select("macaroon").Where("id = ?",
		vendorID.String()).First(&vendor)
	if err != nil {
		return err
	}

	// Generate the delegating Macaroon
	m, err := macaroons.MacaroonFromBytes(vendor.Macaroon.ByteSlice)
	if err != nil {
		return err
	}

	// Add the caveats
	user_id := fmt.Sprintf("user_ID= %s", userID.String())
	delegated, err := vs.AddFirstPartyCaveats(m, []string{user_id})
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
	return c.Render(200, r.HTML("api/vendors/show.html"))
}
