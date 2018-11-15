package macaroons

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop"
	"gopkg.in/macaroon-bakery.v2/bakery"
	"gopkg.in/macaroon-bakery.v2/bakery/checkers"
	"gopkg.in/macaroon-bakery.v2/bakery/postgresrootkeystore"
	macaroon "gopkg.in/macaroon.v2"
)

var dischargeOp = bakery.Op{"firstparty", "x"}
var log logger.FieldLogger
var store bakery.RootKeyStore

// Bakery wraps a bakery.Bakery and provides some nice helper functions
type Bakery struct {
	b        *bakery.Bakery
	oven     *bakery.Oven
	location string
}

type strKey struct{}

func init() {
	log = logger.NewLogger("BAKERY")
	// Create store
	// This is bad, but it seems to work
	db, err := sql.Open("postgres", "host=localhost user=postgres password=postgres database=cms_authz_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Creating root key store")
	ks := postgresrootkeystore.NewRootKeys(db, "root_keys", 100)
	store = ks.NewStore(postgresrootkeystore.Policy{
		ExpiryDuration: 5 * time.Hour,
	})
}

func NewBakery(location string, checker *checkers.Checker, db *pop.Connection) (*Bakery, error) {

	// Do something dumb for public keys
	locator := bakery.NewThirdPartyStore()
	third := bakery.MustGenerateKey()
	locator.AddInfo(location, bakery.ThirdPartyInfo{
		PublicKey: third.Public,
		Version:   bakery.LatestVersion,
	})

	p := bakery.BakeryParams{
		Location:     location,
		Key:          nil,
		Locator:      locator,
		Checker:      checker,
		RootKeyStore: store,
	}

	b := bakery.New(p)

	return &Bakery{
		b:        b,
		oven:     b.Oven,
		location: location,
	}, nil
}

func (b Bakery) NewFirstPartyMacaroon(conditions []string) (*bakery.Macaroon, error) {

	caveats := buildCaveats("", conditions)

	mac, err := b.b.Oven.NewMacaroon(context.Background(), bakery.LatestVersion, caveats, dischargeOp)
	if err != nil {
		return nil, err
	}

	return mac, nil
}

func (b Bakery) AddFirstPartyCaveats(m *bakery.Macaroon, conditions []string) (*bakery.Macaroon, error) {
	caveats := buildCaveats("", conditions)
	err := b.oven.AddCaveats(context.Background(), m, caveats)
	return m, err
}

func (b Bakery) AddThirdPartyCaveat(m *bakery.Macaroon, loc string, conditions []string) (*bakery.Macaroon, error) {

	caveats := buildCaveats(loc, conditions)

	err := b.oven.AddCaveats(context.Background(), m, caveats)
	return m, err
}

func (b Bakery) VerifyMacaroon(ctx context.Context, m *bakery.Macaroon) error {
	ops, conds, err := b.oven.VerifyMacaroon(context.Background(), macaroon.Slice{m.M()})
	if err != nil {
		return err
	}

	fmt.Printf("Ops: %s\n", ops)
	fmt.Printf("Conds: %s\n", conds)

	for _, cond := range conds {
		err := b.b.Checker.CheckFirstPartyCaveat(ctx, cond)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func buildCaveats(location string, conditions []string) []checkers.Caveat {
	caveats := []checkers.Caveat{}

	for _, cond := range conditions {
		caveat := checkers.Caveat{
			Location:  location,
			Condition: cond,
			Namespace: checkers.StdNamespace,
		}
		caveats = append(caveats, caveat)
	}
	return caveats
}

func strContext(key, s string) context.Context {
	return context.WithValue(context.Background(), key, s)
}
