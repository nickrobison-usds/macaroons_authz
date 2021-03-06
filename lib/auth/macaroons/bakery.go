package macaroons

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/gobuffalo/envy"
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
var tstore MemThirdPartyStore
var once sync.Once

// Bakery wraps a bakery.Bakery and provides some nice helper functions
type Bakery struct {
	b        *bakery.Bakery
	oven     *bakery.Oven
	location string
}

func init() {
	log = logger.NewLogger("BAKERY")

	// Get the Database URL from the ENV, or use a default
	url := envy.Get("DATABASE_URL", "")
	if url == "" {
		url = "host=localhost user=raac database=macaroons_authz_development sslmode=disable"
	}
	// Create store
	// This is bad, but it seems to work
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Creating root key store")
	ks := postgresrootkeystore.NewRootKeys(db, "root_keys", 100)
	store = ks.NewStore(postgresrootkeystore.Policy{
		ExpiryDuration: 5 * time.Hour,
	})
	once.Do(func() {
		log.Debug("Creating new Third party key store")
		tstore = NewMemThirdPartyStore()
	})
}

func NewBakery(location string, checker *checkers.Checker, db *pop.Connection, keys *bakery.KeyPair) (*Bakery, error) {

	// Do something dumb for public keys
	if keys == nil {
		log.Debug("Generating pub/priv key pair")
		keys = bakery.MustGenerateKey()
	}
	log.Debugf("Private: %s, Public: %s", keys.Private, keys.Public)
	tstore.AddInfo(location, bakery.ThirdPartyInfo{
		PublicKey: keys.Public,
		Version:   bakery.LatestVersion,
	})

	p := bakery.BakeryParams{
		Logger:       BakedLogger{log},
		Location:     location,
		Key:          keys,
		Locator:      tstore,
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

// GetPrivateKey returns the binary encoding of the Bakery private key
func (b Bakery) GetPrivateKey() []byte {
	key, err := b.oven.Key().Private.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return key
}

// GetPublicKey returns the binary encoding of the Bakery public key
func (b Bakery) GetPublicKey() []byte {
	log.Debug("Getting public key from Bakery: ", []byte(b.oven.Key().Public.String()))
	key, err := b.oven.Key().Public.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return key
}

func (b Bakery) NewFirstPartyMacaroon(conditions []string) (*bakery.Macaroon, error) {

	caveats := buildCaveats("", conditions)

	mac, err := b.b.Oven.NewMacaroon(context.Background(), bakery.LatestVersion, caveats, dischargeOp)
	if err != nil {
		return nil, err
	}

	return mac, nil
}

// NewThirdPartyMacaroon creates a new macaroon with a set of third party caveats, linked to a given locaiton.
func (b Bakery) NewThirdPartyMacaroon(ctx context.Context, loc string, conditions []string) (*bakery.Macaroon, error) {
	caveats := buildCaveats(loc, conditions)

	mac, err := b.oven.NewMacaroon(ctx, bakery.LatestVersion, caveats, dischargeOp)
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

func (b Bakery) VerifyMacaroons(ctx context.Context, m macaroon.Slice) error {

	_, conds, err := b.oven.VerifyMacaroon(context.Background(), m)
	if err != nil {
		return err
	}

	for _, cond := range conds {
		err := b.b.Checker.CheckFirstPartyCaveat(ctx, cond)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b Bakery) VerifyMacaroon(ctx context.Context, m *bakery.Macaroon) error {
	return b.VerifyMacaroons(ctx, macaroon.Slice{m.M()})
}

func (b Bakery) DischargeCaveatByID(ctx context.Context, id string, caveatChecker bakery.ThirdPartyCaveatCheckerFunc) (*bakery.Macaroon, error) {

	// Decode the id
	decodedID, err := macaroon.Base64Decode([]byte(id))
	if err != nil {
		return nil, err
	}

	var caveat []byte

	// Do the discharge thing
	params := bakery.DischargeParams{
		Id:      decodedID,
		Caveat:  caveat,
		Checker: caveatChecker,
		Key:     b.oven.Key(),
		Locator: b.oven.Locator(),
	}

	log.Debug("Discharging it")
	log.Debug("ID: ", string(decodedID))
	log.Debug("ID Bytes: ", decodedID)

	log.Debug("Pub key: ", []byte(b.oven.Key().Public.String()))
	log.Debug("Priv key: ", b.oven.Key().Private.String())

	mac, err := bakery.Discharge(ctx, params)
	if err != nil {
		return nil, err
	}
	return mac, err
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

// BakedLogger wraps a buffalo.Logger so that it works as a bakery.Logger
type BakedLogger struct {
	log logger.FieldLogger
}

// Infof logs info logs
func (b BakedLogger) Infof(_ context.Context, f string, args ...interface{}) {
	b.log.Infof(f, args)
}

// Debugf logs debug logs
func (b BakedLogger) Debugf(_ context.Context, f string, args ...interface{}) {
	b.log.Debugf(f, args)
}
