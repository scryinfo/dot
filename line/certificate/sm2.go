// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tjfoc/gmsm/sm2"

	"github.com/scryinfo/dot/dot"
)

// Sm2 dot
type Sm2 struct {
	logger      *dot.LoggerType
	certificate LoadCertificate
}

func NewRsa(logger *dot.LoggerType) *Sm2 {
	return &Sm2{
		logger:      logger,
		certificate: LoadCertificate{logger: logger},
	}
}

// GenerateRoot Generate ca certificate and private key
// keyFile private key, pemFile ca certificate file
func (c *Sm2) GenerateRoot(rootPri *sm2.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {

	rootCert, err := c.certificate.GenerateRoot(x509.PureEd25519, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateRootFile(rootPri, rootCert, rootPri.Public(), keyFile, pemFile)

	return rootCert, err
}

// GenerateLeaf Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *Sm2) GenerateLeaf(rootCert *x509.Certificate, rootPri *sm2.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {
	leaf, err := c.certificate.GenerateLeafCertificate(x509.PureEd25519, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	leafPri, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateLeafFile(leafPri, leaf, leafPri.Public(), keyFile, pemFile, rootCert, rootPri)
	return leaf, err
}

// PrivateKey Read private key from keyFile
func (c *Sm2) PrivateKey(keyFile string) (*sm2.PrivateKey, error) {

	key, err := c.certificate.LoadPrivateKey(keyFile)
	if err != nil {
		return nil, err
	}

	if resKey, ok := key.(*sm2.PrivateKey); ok {
		return resKey, nil
	} else {
		return nil, fmt.Errorf("the key isnt sm2.PrivateKey")
	}
}

// PublicKey Read public key from certFile
func (c *Sm2) PublicKey(certFile string) (*sm2.PublicKey, error) {
	cer, err := c.certificate.LoadCertificate(certFile)
	if err != nil {
		return nil, err
	}

	if tpub, ok := cer.PublicKey.(*sm2.PublicKey); ok {
		return tpub, nil
	} else {
		return nil, errors.New("do not sm2.PublicKey")
	}
}

// MakeSm2Key Generate private key
func MakeSm2Key() (*sm2.PrivateKey, error) {
	key, err := sm2.GenerateKey(rand.Reader)
	return key, err
}
