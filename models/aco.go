package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// Data model representing an ACO
type ACO struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Macaroon    nulls.ByteSlice `json:"macaroon" db:"macaroon"`
	Certificate Certificate     `json:"certificates" has_many:"certificates" fd_id:"aco_id"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// Returns the appropriate TableName, which addresses the weird pluralization
func (a ACO) TableName() string {
	return "acos"
}

// Returns the ACO ID field, as a string
func (a ACO) StringID() string {
	return a.ID.String()
}

// String is not required by pop and may be deleted
func (a ACO) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// ACOS is not required by pop and may be deleted
type ACOS []ACO

// String is not required by pop and may be deleted
func (a ACOS) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

func (a *ACO) AfterCreate(tx *pop.Connection) error {
	log.Debugf("Loading %s and %s\n", a.ID, a.Certificate.ACOID)
	return tx.Save(&a.Certificate)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *ACO) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *ACO) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *ACO) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
