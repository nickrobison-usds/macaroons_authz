package grifts

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/markbates/grift/grift"
	"github.com/nickrobison-usds/macaroons_authz/actions"
	"github.com/nickrobison-usds/macaroons_authz/lib/helpers"
	"github.com/nickrobison-usds/macaroons_authz/models"
	"github.com/pkg/errors"
)

var log logger.FieldLogger

var _ = grift.Namespace("db", setupActionsDatabase)

func init() {
	log = logger.NewLogger("DB:SEED")
}

func setupActionsDatabase() {
	// DB seeding
	grift.Desc("seed", "Seed the database with some initial data")
	grift.Add("seed", seedDatabase)
}

func seedDatabase(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		log.Debug("Truncating database tables")
		if err := tx.TruncateAll(); err != nil {
			return errors.WithStack(err)
		}

		// ACOS
		log.Debug("Loading ACO seeds")
		areader, err := getCSVReader("./db/seeds_aco.csv")
		if err != nil {
			return errors.WithStack(err)
		}

		err = processCSV(areader, tx, deserializeACO)
		if err != nil {
			return errors.WithStack(err)
		}

		// Users
		log.Debug("Loading User seeds")
		ureader, err := getCSVReader("./db/seeds_user.csv")
		if err != nil {
			return errors.WithStack(err)
		}

		err = processCSV(ureader, tx, deserializeUser)
		if err != nil {
			return errors.WithStack(err)
		}

		// Vendors
		log.Debug("Loading Vendor seeds")
		vreader, err := getCSVReader("./db/seeds_vendor.csv")
		if err != nil {
			return errors.WithStack(err)
		}

		err = processCSV(vreader, tx, deserializeVendor)
		if err != nil {
			return errors.WithStack(err)
		}

		log.Debug("Seeding finished")
		return nil
	})

}

func getCSVReader(filename string) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Get the CSV reader
	reader := csv.NewReader(file)

	return reader, nil
}

func processCSV(reader *csv.Reader, tx *pop.Connection, deserializer func(record []string, tx *pop.Connection) error) error {
	// Read each of the rows and process it
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}

		err = deserializer(row, tx)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func deserializeACO(record []string, tx *pop.Connection) error {
	aco := models.ACO{}

	aco.ID = helpers.MustGenerateID()
	aco.Name = record[0]

	// Add the certs and macaroons
	err := actions.CreateACOCertificates(&aco)
	if err != nil {
		return err
	}

	return tx.Create(&aco)
}

func deserializeUser(record []string, tx *pop.Connection) error {
	user := models.User{}
	user.ID = helpers.MustGenerateID()

	user.Name = record[0]
	user.Email = nulls.NewString(record[1])
	user.Provider = record[2]
	user.ProviderID = record[3]

	return tx.Create(&user)
}

func deserializeVendor(record []string, tx *pop.Connection) error {
	vendor := models.Vendor{}
	vendor.ID = helpers.MustGenerateID()
	err := actions.CreateVendorCertificates(&vendor)
	if err != nil {
		return err
	}

	vendor.Name = record[0]

	return tx.Create(&vendor)
}
