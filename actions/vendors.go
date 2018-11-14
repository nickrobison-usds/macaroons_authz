package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
)

// VendorsCreate default implementation.
func VendorsCreate(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/vendors/create.html"))
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

// VendorsShow default implementation.
func VendorsShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/vendors/show.html"))
}
