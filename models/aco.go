package models

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/nickrobison/cms_authz/lib/auth/ca"
	"github.com/pkg/errors"
	macaroon "gopkg.in/macaroon.v2"
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
	return tx.Save(&a.Certificate)
}

func (a *ACO) BeforeCreate(tx *pop.Connection) error {
	// Set the UUID
	id := mustGenerateUUID()
	a.ID = id

	// Now do the cert thing
	log.Debug("Creating CA")

	cert, err := ca.CreateCA(a.Name, "aco")
	if err != nil {
		return err
	}

	fmt.Println(cert)

	parsed, err := ca.ParseCFSSLResponse(&cert)
	if err != nil {
		return errors.WithStack(err)
	}

	encCert, err := parsed.EncodeCertificate()
	if err != nil {
		return errors.WithStack(err)
	}

	priv, err := parsed.EncodePrivateKey()
	if err != nil {
		return errors.WithStack(err)
	}

	acoCert := Certificate{
		ACOID:       id,
		Key:         priv,
		Certificate: encCert,
		SHA:         parsed.SHA,
	}

	a.Certificate = acoCert

	// Now, generate the macaroon

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	nonceString := base64.StdEncoding.EncodeToString(nonce)

	// We need to use V1 serialization, for some reason.
	log.Debug(a.Certificate.Key)
	m, err := macaroon.New([]byte("test key"), []byte(nonceString), "http://localhost:8080", macaroon.V1)
	if err != nil {
		return err
	}

	// Add the first party claims

	caveats := map[string]string{
		"aco_id": id.String(),
	}
	for cav := range caveats {
		err = m.AddFirstPartyCaveat([]byte(cav))
		if err != nil {
			return err
		}
	}

	mBinary, err := m.MarshalBinary()
	if err != nil {
		return err
	}
	a.Macaroon = nulls.NewByteSlice(mBinary)

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
