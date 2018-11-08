package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/nickrobison/cms_authz/lib/auth/ca"
)

// Data model representing an ACO
type ACO struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Certificate Certificate `json:"certificates" has_many:"certificates" fd_id:"aco_id"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
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

func (a *ACO) BeforeCreate(tx *pop.Connection) error {
	// Set the UUID
	a.ID = mustGenerateUUID()

	// Now do the cert thing
	log.Debug("Creating CA")

	cert, err := ca.CreateCA(a.Name, "aco")
	if err != nil {
		return err
	}

	a.Certificate.Key = cert.Certificate
	a.Certificate.Certificate = cert.Certificate
	a.Certificate.SHA = cert.Sums.Certificate.SHA1

	return nil
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
