package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// Data model representing an ACO
// Storing the Private Key and cert here is NOT A GOOD IDEA! It's only temporary.
type ACO struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Key         string    `json:"key"  db:"key"`
	Certificate string    `json:"certificate" db:"certificate"`
	SHA         string    `json:"sha" db:"sha"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Returns the appropriate TableName, which addresses the weird pluralization
func (a ACO) TableName() string {
	return "acos"
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
