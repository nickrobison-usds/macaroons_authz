package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison/cms_authz/lib/auth/ca"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
)

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

func RenderCreatePage(c buffalo.Context) error {
	aco := models.ACO{}

	c.Set("aco", aco)
	return c.Render(http.StatusOK, r.HTML("api/acos/create.html"))
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

	// Generate a new ID for the ACO
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.WithStack(err)
	}
	aco.ID = uid

	err = c.Bind(&aco)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("\n\n\nACO: %v\n\n\n", aco)

	// Try to create a new CA
	cert, err := ca.CreateCA(aco.Name, "aco")
	if err != nil {
		return errors.WithStack(err)
	}

	aco.Key = cert.Certificate
	aco.Certificate = cert.Certificate
	aco.SHA = cert.Sums.Certificate.SHA1

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Create(&aco); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/api/acos/index")
}
