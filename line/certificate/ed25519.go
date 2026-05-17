// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"

	"github.com/scryinfo/dot/dot"
)

// Ed25519 dot
type Ed25519 struct {
	logger      *dot.LoggerType
	certificate BaseCertificate
}

func NewEd25519(logger *dot.LoggerType) *Ed25519 {
	return &Ed25519{
		logger:      logger,
		certificate: BaseCertificate{logger: logger},
	}
}

// GenerateRoot Generate ca certificate and private key
// keyFile private key, pemFile ca certificate file
func (c *Ed25519) GenerateRoot(rootPri ed25519.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {

	rootCert, err := c.certificate.GenerateRoot(x509.UnknownSignatureAlgorithm, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateRootFile(rootPri, rootCert, rootPri.Public(), keyFile, pemFile)

	return rootCert, err
}

// GenerateLeaf Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *Ed25519) GenerateLeaf(rootCert *x509.Certificate, rootPri ed25519.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {
	leaf, err := c.certificate.GenerateLeafCertificate(x509.UnknownSignatureAlgorithm, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	_, leafPri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateLeafFile(leafPri, leaf, leafPri.Public(), keyFile, pemFile, rootCert, rootPri)
	return leaf, err
}

// PrivateKey Read private key from keyFile
func (c *Ed25519) PrivateKey(keyFile string) (*ed25519.PrivateKey, error) {

	key, err := c.certificate.LoadPrivateKey(keyFile)
	if err != nil {
		return nil, err
	}

	if resKey, ok := key.(*ed25519.PrivateKey); ok {
		return resKey, nil
	} else {
		return nil, fmt.Errorf("the key isnt ed25519.PrivateKey")
	}
}

// PublicKey Read public key from certFile
func (c *Ed25519) PublicKey(certFile string) (*ed25519.PublicKey, error) {
	cer, err := c.certificate.LoadCertificate(certFile)
	if err != nil {
		return nil, err
	}

	if tpub, ok := cer.PublicKey.(*ed25519.PublicKey); ok {
		return tpub, nil
	} else {
		return nil, errors.New("do not ed25519.PublicKey")
	}
}

// MakeEd25519Key Generate private key
func MakeEd25519Key() (ed25519.PrivateKey, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	return key, err
}
