package actions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison-usds/macaroons_authz/lib/auth/macaroons"
	"github.com/nickrobison-usds/macaroons_authz/lib/helpers"
	"github.com/nickrobison-usds/macaroons_authz/models"
	"github.com/pkg/errors"
	"github.com/rakutentech/jwk-go/jwk"
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

// dischargeResponse contains the response from a /discharge POST request.
type dischargeResponse struct {
	Macaroon *bakery.Macaroon `json:",omitempty"`
}

type Identifiable interface {
	GetID() uuid.UUID
	TableName() string
}

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
	s, err := macaroons.NewBakery(userURI, checkers.New(nil), models.DB, key)
	if err != nil {
		log.Fatal(err)
	}
	us = s
}

// UsersFind looks up a userid for a given user name
func UsersFind(c buffalo.Context) error {
	nameString := c.Param("name")
	if nameString == "" {
		return c.Render(http.StatusBadRequest, r.String("Cannot have blank query name."))
	}

	user := models.User{}

	tx := c.Value("tx").(*pop.Connection)
	err := tx.Where("name = ?", nameString).First(&user)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, r.String(user.ID.String()))
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

// UsersTokenGet retrieves the entity token associated with the given user
func UsersTokenGet(c buffalo.Context) error {
	userID := c.Param("user_id")
	entity_type := c.Param("entity_type")
	entity_id := c.Param("entity_id")

	tx := c.Value("tx").(*pop.Connection)

	var mac_bytes []byte

	// Get the model type
	switch entity_type {
	case "ACO":
		{
			user := models.AcoUser{}
			err := tx.Select("macaroon").Where("aco_id = ? AND entity_id = ?", entity_id, userID).First(&user)
			if err != nil {
				return errors.WithStack(err)
			}
			mac_bytes = user.Macaroon
		}
	case "Vendor":
		{
			user := models.VendorUser{}
			err := tx.Select("macaroon").Where("vendor_id = ? AND user_id = ?").First(&user)
			if err != nil {
				return errors.WithStack(err)
			}
			mac_bytes = user.Macaroon

		}
	default:
		return errors.WithStack(fmt.Errorf("Cannot get token for entity type %s", entity_type))
	}

	// return it
	return c.Render(200, r.String(macaroons.EncodeMacaroon(mac_bytes)))
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
			err := AssignUserToACO(helpers.UUIDOfString(req.EntityID), helpers.UUIDOfString(req.UserID), tx)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	case "Vendor":
		{
			err := AssignUserToVendor(helpers.UUIDOfString(req.EntityID), helpers.UUIDOfString(req.UserID), tx)
			if err != nil {
				return errors.WithStack(err)
			}
		}

	default:
		return errors.WithStack(errors.New("Cannot create a non-User type"))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/api/users/index")
}

// UsersJWKS returns the jwks.json for the given third-party
func UsersJWKS(c buffalo.Context) error {
	pub := us.GetPublicKey()
	log.Debug("Public key from bakery:", string(pub))
	spec := jwk.NewSpec(pub)

	spec.KeyID = "1"
	spec.Algorithm = "ES256"
	spec.Use = "enc"

	log.Debug("Spec: ", spec)

	json, err := spec.MarshalJSON()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Debug(string(json))
	return c.Render(http.StatusOK, r.JSON(spec))
}

func UsersVerify(c buffalo.Context) error {
	token := c.Param("id64")

	// If it comes in as id, we need to translate it to base64 encoding
	if token == "" {
		token = base64.URLEncoding.EncodeToString(byte(c.Param("id")))
	}

	log.Debugf("Token: %s", token)

	log.Debug("Discharging")

	mac, err := us.DischargeCaveatByID(c.Request().Context(), token, userIDCaveatChecker(c.Value("tx").(*pop.Connection)))
	if err != nil {
		return errors.WithStack(err)
	}

	b, err := mac.MarshalJSON()
	if err != nil {
		return errors.WithStack(err)
	}

	log.Debugf("Marshalled: %s", string(b))

	return c.Render(200, r.JSON(&dischargeResponse{mac}))
}

func userIDCaveatChecker(db *pop.Connection) bakery.ThirdPartyCaveatCheckerFunc {

	return func(ctx context.Context, cav *bakery.ThirdPartyCaveatInfo) ([]checkers.Caveat, error) {

		var caveats []checkers.Caveat
		log.Debug("In the ID checker")
		log.Debug(ctx)
		log.Debug(string(cav.Condition))
		_, arg, err := checkers.ParseCaveat(string(cav.Condition))
		if err != nil {
			return caveats, err
		}
		log.Debug("UUID: ", arg)

		var user models.User

		// Getting from the DB
		ID := helpers.UUIDOfString(arg)

		err = db.Select("id").Where("id = ? ", ID).First(&user)
		if err != nil {
			return nil, err
		}

		if user.ID != ID {
			return nil, bakery.ErrPermissionDenied
		}

		return nil, nil
	}

}
