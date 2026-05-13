// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/pkg/errors"

	"github.com/scryinfo/dot/dot"
)

const (
	//EcdsaTypeID type id of dot
	EcdsaTypeID = "4b8b1751-4799-4578-af46-d9b339cf582f"
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

	rootCa, err := c.certificate.GenerateRoot(x509.ECDSAWithSHA256, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	err = c.certificate.GenerateRootFile(rootPri, rootCa, rootPri.Public(), keyFile, pemFile)

	return rootCa, err
}

// GenerateLeaf Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *Ecdsa) GenerateLeaf(rootCa *x509.Certificate, rootPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (*x509.Certificate, error) {
	leaf, err := c.certificate.GenerateLeafCertificate(x509.ECDSAWithSHA256, dnsName, orgName)
	if err != nil {
		return nil, err
	}
	leafPri, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	err = c.certificate.GenerateLeafFile(leafPri, leaf, leafPri.Public(), keyFile, pemFile, rootCa, rootPri)
	return leaf, err
}

// PrivateKey Read private key from keyFile
func (c *Ecdsa) PrivateKey(keyFile string) (pri *ecdsa.PrivateKey, err error) {

	file, err := wdFile(keyFile)

	if err != nil || len(file) < 1 {
		return nil, errors.New("file do not exist")
	}

	bs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bs)
	if block == nil || block.Bytes == nil {
		return nil, errors.New("do not parse the data of file")
	}

	pri, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pri, err
}

// PublicKey Read public key from pemFile
func (c *Ecdsa) PublicKey(pemFile string) (pubKey *ecdsa.PublicKey, err error) {

	cer, err := c.Certificate(pemFile)
	if err != nil {
		return nil, err
	}

	tpub, ok := cer.PublicKey.(*ecdsa.PublicKey)
	if ok {
		pubKey = tpub
	} else {
		return nil, errors.New("do not ecdsa.PublicKey")
	}
	return pubKey, err
}

// MakePriKey Generate private key
func MakePriKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
