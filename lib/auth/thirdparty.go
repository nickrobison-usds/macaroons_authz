package auth

import (
	"context"

	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop"
	"gopkg.in/macaroon-bakery.v2/bakery"
)

// KeyPair masks the underlying bakery.KeyPair
type KeyPair bakery.KeyPair

// Interface for retrieving a Public/Private KeyPair from an application model
type KeyPairGetter interface {
	KeyPair() KeyPair
}

type ApplicationThirdPartyStore struct {
	db    *pop.Connection
	model KeyPairGetter
}

var log logger.FieldLogger

func init() {
	log = logger.NewLogger("AUTH")
}

func (l ApplicationThirdPartyStore) ThirdPartyInfo(ctx context.Context, loc string) (bakery.ThirdPartyInfo, error) {
	var info bakery.ThirdPartyInfo

	err := l.db.First(&l.model)
	if err != nil {
		return info, err
	}

	log.Debug(l.model.KeyPair())

	return info, nil
}
