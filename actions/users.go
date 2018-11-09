package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
)

type userAssignRequest struct {
	UserID     string `form:"userID"`
	EntityType string `form:"assignEntity"`
	EntityID   string `form:"entityOptions"`
}

// UsersIndex default implementation.
func UsersIndex(c buffalo.Context) error {

	users := []models.User{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.All(&users); err != nil {
		return errors.WithStack(err)
	}
	c.Set("users", users)

	return c.Render(200, r.HTML("api/users/index.html"))
}

// UsersShow default implementation.
func UsersShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/users/show.html"))
}

// UsersCreate default implementation.
func UsersCreate(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/users/create.html"))
}

// UsersDelete default implementation.
func UsersDelete(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/users/delete.html"))
}

// UsersAssign adds a given user as an authorized member of the entity
func UsersAssign(c buffalo.Context) error {
	req := userAssignRequest{}

	err := c.Bind(&req)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Printf("Assigning: %s\n", req.EntityID)
	tx := c.Value("tx").(*pop.Connection)

	switch req.EntityType {
	case "ACO":
		{
			err := delegateUserToACO(req, tx)
			if err != nil {
				return errors.WithStack(err)
			}
		}

	default:
		return errors.WithStack(errors.New("Cannot create a non-User type"))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/api/users/index")
}

func delegateUserToACO(req userAssignRequest, tx *pop.Connection) error {
	link := models.AcoUser{
		ACOID:  UUIDOfString(req.EntityID),
		UserID: UUIDOfString(req.UserID),
	}

	return tx.Save(&link)
}

// UUIDOfString always returns a uuid.UUID, otherwise it panics
func UUIDOfString(id string) uuid.UUID {
	str, err := uuid.FromString(id)
	if err != nil {
		panic(err)
	}
	return str
}
