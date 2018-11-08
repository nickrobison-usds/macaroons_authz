package models

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection
var log logger.FieldLogger

func init() {
	log = logger.NewLogger("MODELS")
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}

// Forces the generation of a UUID, panics if it can't.
func mustGenerateUUID() uuid.UUID {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return uuid
}
