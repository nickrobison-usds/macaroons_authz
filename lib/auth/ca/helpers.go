package ca

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

// PublicPrivateKeyPair contains all the necessary cryptographic data for doing useful things.
type PublicPrivateKeyPair struct {
	PublicKey   *ecdsa.PublicKey
	PrivateKey  *ecdsa.PrivateKey
	Certificate *x509.Certificate
	SHA         string
}

// EncodePublicKey marshalls the PublicKey as a string
func (p PublicPrivateKeyPair) EncodePublicKey() (string, error) {
	mBytes, err := x509.MarshalPKIXPublicKey(p.PublicKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(mBytes), nil
}

// EncodePrivateKey marshalls the private key as a string
func (p PublicPrivateKeyPair) EncodePrivateKey() (string, error) {
	mBytes, err := x509.MarshalECPrivateKey(p.PrivateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(mBytes), nil
}

func (p PublicPrivateKeyPair) EncodeCertificate() (string, error) {
	return hex.EncodeToString(p.Certificate.Raw), nil
}

// ParseCFSSLResponse takes a CFSSLCertificateResponse and parses into something that uses the standard go primitives.
func ParseCFSSLResponse(resp *CFSSLCertificateResponse) (*PublicPrivateKeyPair, error) {
	pair := &PublicPrivateKeyPair{
		SHA: resp.Sums.Certificate.SHA1,
	}

	privPem, _ := pem.Decode([]byte(resp.PrivateKey))

	priv, err := x509.ParseECPrivateKey(privPem.Bytes)
	if err != nil {
		return pair, err
	}

	certPem, _ := pem.Decode([]byte(resp.Certificate))

	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return pair, err
	}

	// Get the public key (and assert it)
	pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return pair, errors.New("Cannot cast PublicKey to ECDSA")
	}
	pair.PrivateKey = priv
	pair.PublicKey = pub
	pair.Certificate = cert

	return pair, nil
}
