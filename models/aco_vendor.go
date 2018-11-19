package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type AcoVendor struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ACOID     uuid.UUID `json:"aco_id" db:"aco_id"`
	VendorID  uuid.UUID `json:"vendor_id" db:"vendor_id"`
	Macaroon  []byte    `json:"macaroon" db:"macaroon"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a AcoVendor) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AcoVendors is not required by pop and may be deleted
type AcoVendors []AcoVendor

// String is not required by pop and may be deleted
func (a AcoVendors) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AcoVendor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AcoVendor) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AcoVendor) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
