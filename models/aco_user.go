package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/nickrobison/cms_authz/lib/auth/macaroons"
)

type AcoUser struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ACOID     uuid.UUID `json:"aco" db:"aco_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Macaroon  []byte    `json:"macaroon" db:"macaroon"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a AcoUser) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AcoUsers is not required by pop and may be deleted
type AcoUsers []AcoUser

// String is not required by pop and may be deleted
func (a AcoUsers) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AcoUser) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AcoUser) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AcoUser) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// BeforeCreate generates a delegated macaroon from the ACO.
func (a *AcoUser) BeforeCreate(tx *pop.Connection) error {
	// Get the Macaroon from the ACO
	aco := ACO{}

	err := tx.Select("macaroon").Where("id = ?",
		a.ACOID.String()).First(&aco)
	if err != nil {
		return err
	}

	// Generate a delegating Macaroon
	m, err := macaroons.MacaroonFromBytes(aco.Macaroon.ByteSlice)
	if err != nil {
		return err
	}

	delegated, err := macaroons.DelegateACOToUser(a.ACOID, a.UserID, &m)
	if err != nil {
		return err
	}
	mBinary, err := delegated.MarshalBinary()
	if err != nil {
		return err
	}

	a.Macaroon = mBinary
	return nil
}
