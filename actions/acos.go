package actions

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
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
	c.Set("binary", showBytes)

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

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Create(&aco); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/api/acos/index")
}

func showBytes(s nulls.ByteSlice) string {
	log.Debugf("In the renderer: %s\n", s.ByteSlice)
	return hex.EncodeToString(s.ByteSlice)
}
