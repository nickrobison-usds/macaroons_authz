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

	acos := []models.Aco{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.All(&acos); err != nil {
		return errors.WithStack(err)
	}

	c.Set("acos", acos)
	return c.Render(200, r.HTML("api/acos/index.html"))
}

func AcosCreate(c buffalo.Context) error {
	aco := models.Aco{}
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.WithStack(err)
	}
	aco.ID = uid

	c.Set("aco", aco)
	return c.Render(http.StatusOK, r.HTML("api/acos/create.html"))
}

func AcosCreateACO(c buffalo.Context) error {
	fmt.Println(c.Request())

	aco := models.Aco{}

	c.Bind(&aco)

	fmt.Printf("\n\n\nACO: %v\n\n\n", aco)

	// Try to create a new CA
	err := ca.CreateCA(aco.Name, "aco")
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	if err := tx.Create(&aco); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/api/acos/index")
}
