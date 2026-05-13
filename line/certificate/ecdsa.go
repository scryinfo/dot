// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"

	"github.com/scryinfo/dot/dot"
)

// Ecdsa dot
type Ecdsa struct {
	logger      *dot.LoggerType
	certificate LoadCertificate
}

func NewEcdsa(logger *dot.LoggerType) *Ecdsa {
	return &Ecdsa{
		logger:      logger,
		certificate: LoadCertificate{logger: logger},
	}
}

// GenerateRoot Generate ca certificate and private key
// keyFile private key, pemFile ca certificate file
func (c *Ecdsa) GenerateRoot(rootPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {

	rootCert, err := c.certificate.GenerateRoot(x509.ECDSAWithSHA256, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateRootFile(rootPri, rootCert, rootPri.Public(), keyFile, pemFile)

	return rootCert, err
}

// GenerateLeaf Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *Ecdsa) GenerateLeaf(rootCert *x509.Certificate, rootPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {
	leaf, err := c.certificate.GenerateLeafCertificate(x509.ECDSAWithSHA256, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	leafPri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateLeafFile(leafPri, leaf, leafPri.Public(), keyFile, pemFile, rootCert, rootPri)
	return leaf, err
}

// PrivateKey Read private key from keyFile
func (c *Ecdsa) PrivateKey(keyFile string) (*ecdsa.PrivateKey, error) {

	key, err := c.certificate.LoadPrivateKey(keyFile)
	if err != nil {
		return nil, err
	}

	if ecdsaKey, ok := key.(*ecdsa.PrivateKey); ok {
		return ecdsaKey, nil
	} else {
		return nil, fmt.Errorf("the key isnt ecdsa.PrivateKey")
	}
}

// PublicKey Read public key from certFile
func (c *Ecdsa) PublicKey(certFile string) (*ecdsa.PublicKey, error) {
	cer, err := c.certificate.LoadCertificate(certFile)
	if err != nil {
		return nil, err
	}

	if tpub, ok := cer.PublicKey.(*ecdsa.PublicKey); ok {
		return tpub, nil
	} else {
		return nil, errors.New("do not ecdsa.PublicKey")
	}
}

// MakeECDSAKey Generate private key
func MakeECDSAKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
