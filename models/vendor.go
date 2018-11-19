package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Vendor struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	Name      string      `json:"name" db:"name"`
	Macaroon  []byte      `json:"macaroon" db:"macaroon"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	Users     VendorUsers `json:"vendor_users" has_many:"vendor_users"`
	ACOTokens AcoVendors  `json:"aco_tokens" has_many:"aco_vendors"`
}

// String is not required by pop and may be deleted
func (v Vendor) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// StringID returns the uuid.UUID as a String
func (v Vendor) StringID() string {
	return v.ID.String()
}

// Vendors is not required by pop and may be deleted
type Vendors []Vendor

// String is not required by pop and may be deleted
func (v Vendors) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (v *Vendor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (v *Vendor) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (v *Vendor) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
