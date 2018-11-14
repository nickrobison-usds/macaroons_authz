package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/nickrobison/cms_authz/lib/helpers"
)

type Certificate struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ACOID       uuid.UUID `json:"aco" db:"aco_id"`
	Key         string    `json:"private_key" db:"key"`
	Certificate string    `json:"certificate" db:"certificate"`
	SHA         string    `json:"sha_sum" db:"sha_sum"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (c Certificate) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Certificates is not required by pop and may be deleted
type Certificates []Certificate

// String is not required by pop and may be deleted
func (c Certificates) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

func (c *Certificate) BeforeCreate(tx *pop.Connection) error {
	c.ID = helpers.MustGenerateID()
	return nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Certificate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Certificate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Certificate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
