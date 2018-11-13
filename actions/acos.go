package actions

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
	"github.com/nickrobison/cms_authz/models"
	"github.com/pkg/errors"
	"gopkg.in/macaroon-bakery.v1/bakery"
	macaroon "gopkg.in/macaroon.v1"
)

type idNamePair struct {
	ID   string
	Name string
}

var service *bakery.Service

func init() {
	p := bakery.NewServiceParams{
		Location: "http://test.loc",
		Store:    nil,
		Key:      nil,
		Locator:  nil,
	}

	b, err := bakery.NewService(p)
	if err != nil {
		log.Fatal(err)
	}
	service = b
}

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

func AcosHeadIndex(c buffalo.Context) error {
	acos := []models.ACO{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Select("id, name").All(&acos); err != nil {
		return errors.WithStack(err)
	}

	values := []idNamePair{}

	for _, aco := range acos {
		values = append(values, idNamePair{Name: aco.Name, ID: aco.StringID()})
	}

	return c.Render(200, r.JSON(&values))

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
	c.Set("bytesToString", showBytes)

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

func AcoVerifyUser(c buffalo.Context) error {
	var requestData models.AcoUser
	err := c.Bind(&requestData)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something bad happened."))
	}

	log.Debugf("Verifying that user %s is a member of %s", requestData.UserID, requestData.ACOID)

	// Check that the association actually exists.
	tx := c.Value("tx").(*pop.Connection)

	var acoUser models.AcoUser

	err = tx.Where("aco_id = ?", requestData.ACOID).Where("user_id = ?", requestData.UserID).First(&acoUser)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something wen't wrong: %s", err.Error()))
	}

	// If it exists, discharge it
	reqM, err := macaroons.MacaroonFromBytes(requestData.Macaroon)
	if err != nil {
		log.Error(err)
		return c.Render(http.StatusInternalServerError, r.String("Something wen't wrong: %s", err.Error()))
	}

	d, err := bakery.DischargeAll(&reqM, dischargeUserCaveat)
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.String("Something went wrong: %s", err.Error()))
	}

	// Inspect everything
	for i, mac := range d {
		log.Debug("Macaroon: ", i)
		log.Debug(mac)
	}

	return c.Render(http.StatusOK, r.String("ok"))
}

// Returns a byte array as a string
func showBytes(s nulls.ByteSlice) string {
	return base64.URLEncoding.EncodeToString(s.ByteSlice)
}

func dischargeUserCaveat(firstPartyLocation string, cav macaroon.Caveat) (*macaroon.Macaroon, error) {

	log.Debug("First party: ", firstPartyLocation)
	log.Debug(cav.Id)

	mac, err := service.Discharge(macaroons.StrcmpChecker("user_id = test"), cav.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return mac, nil
}
