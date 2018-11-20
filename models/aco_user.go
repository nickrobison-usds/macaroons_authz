package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type AcoUser struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ACOID     uuid.UUID `json:"aco_id" db:"aco_id"`
	EntityID  uuid.UUID `json:"entity_id" db:"entity_id"`
	Macaroon  []byte    `json:"macaroon" db:"macaroon"`
	IsUser    bool      `json:"is_user" db:"is_user"`
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
