package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"github.com/nickrobison/cms_authz/lib/helpers"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
)

var userURI = "http://localhost:8080/api/users/verify"

var us *macaroons.Bakery

type userAssignRequest struct {
	UserID     string `form:"userID"`
	EntityType string `form:"assignEntity"`
	EntityID   string `form:"entityOptions"`
}

func init() {
	s, err := macaroons.NewBakery(userURI, checkers.New(nil), models.DB)
	if err != nil {
		log.Fatal(err)
	}
	us = s
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
	userID := c.Param("id")
	user := models.User{}

	tx := c.Value("tx").(*pop.Connection)
	err := tx.Find(&user, userID)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Load(&user)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("user", user)
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
			err := DelegateACOToUser(helpers.UUIDOfString(req.EntityID), helpers.UUIDOfString(req.UserID), tx)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	case "Vendor":
		{
			err := DelegateUserToVendor(helpers.UUIDOfString(req.EntityID), helpers.UUIDOfString(req.UserID), tx)
			if err != nil {
				return errors.WithStack(err)
			}
		}

	default:
		return errors.WithStack(errors.New("Cannot create a non-User type"))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/api/users/index")
}

func UsersVerify(c buffalo.Context) error {
	token := c.Param("id64")

	log.Debugf("Token: %s", token)

	log.Debug("Discharging")

	_, err := us.DischargeCaveatByID(c.Request().Context(), token, caveatChecker)
	if err != nil {
		return errors.WithStack(err)
	}

	log.Debug("Verification complete.")
	return c.Render(http.StatusUnauthorized, r.String("Nope"))
}

func caveatChecker(ctx context.Context, cav *bakery.ThirdPartyCaveatInfo) ([]checkers.Caveat, error) {
	log.Debug("In the checker")
	log.Debug(ctx)
	log.Debug(cav)
	return []checkers.Caveat{}, nil
}
