package macaroons

import (
	"context"

	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
)

var dischargeOp = bakery.Op{"firstparty", "x"}

type Bakery struct {
	b        *bakery.Bakery
	location string
}

func NewBakery(location string) (*Bakery, error) {

	// Do something dumb for public keys
	locator := bakery.NewThirdPartyStore()
	third := bakery.MustGenerateKey()
	locator.AddInfo(location, bakery.ThirdPartyInfo{
		PublicKey: third.Public,
		Version:   bakery.LatestVersion,
	})

	p := bakery.BakeryParams{
		Location: location,
		Key:      nil,
		Locator:  locator,
	}

	return &Bakery{
		b:        bakery.New(p),
		location: location,
	}, nil
}

func (b Bakery) NewFirstPartyMacaroon(conditions []string) (*bakery.Macaroon, error) {
	caveats := []checkers.Caveat{}

	for _, cond := range conditions {
		caveat := checkers.Caveat{
			Condition: cond,
		}
		caveats = append(caveats, caveat)
	}
	return b.b.Oven.NewMacaroon(context.Background(), bakery.LatestVersion, caveats, dischargeOp)
}

func AddThirdPartyCaveat(m *bakery.Macaroon, loc string, conditions []string) (*bakery.Macaroon, error) {
	caveats := []checkers.Caveat{}

	for _, cond := range conditions {
		caveat := checkers.Caveat{
			Location:  loc,
			Condition: cond,
		}
		caveats = append(caveats, caveat)
	}

	err := m.AddCaveats(context.Background(), caveats, nil, nil)
	return m, err
}
