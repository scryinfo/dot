// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"

	"github.com/scryinfo/dot/dot"
)

const rsaKeySize = 2048

// Rsa dot
type Rsa struct {
	logger      *dot.LoggerType
	certificate LoadCertificate
}

func NewSm2(logger *dot.LoggerType) *Rsa {
	return &Rsa{
		logger:      logger,
		certificate: LoadCertificate{logger: logger},
	}
}

// GenerateRoot Generate ca certificate and private key
// keyFile private key, pemFile ca certificate file
func (c *Rsa) GenerateRoot(rootPri *rsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {

	rootCert, err := c.certificate.GenerateRoot(x509.PureEd25519, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateRootFile(rootPri, rootCert, rootPri.Public(), keyFile, pemFile)

	return rootCert, err
}

// GenerateLeaf Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *Rsa) GenerateLeaf(rootCert *x509.Certificate, rootPri *rsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {
	leaf, err := c.certificate.GenerateLeafCertificate(x509.PureEd25519, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	leafPri, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateLeafFile(leafPri, leaf, leafPri.Public(), keyFile, pemFile, rootCert, rootPri)
	return leaf, err
}

// PrivateKey Read private key from keyFile
func (c *Rsa) PrivateKey(keyFile string) (*rsa.PrivateKey, error) {

	key, err := c.certificate.LoadPrivateKey(keyFile)
	if err != nil {
		return nil, err
	}

	if resKey, ok := key.(*rsa.PrivateKey); ok {
		return resKey, nil
	} else {
		return nil, fmt.Errorf("the key isnt rsa.PrivateKey")
	}
}

// PublicKey Read public key from certFile
func (c *Rsa) PublicKey(certFile string) (*rsa.PublicKey, error) {
	cer, err := c.certificate.LoadCertificate(certFile)
	if err != nil {
		return nil, err
	}

	if tpub, ok := cer.PublicKey.(*rsa.PublicKey); ok {
		return tpub, nil
	} else {
		return nil, errors.New("do not rsa.PublicKey")
	}
}

// MakeRsaKey Generate private key
func MakeRsaKey() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	return key, err
}
