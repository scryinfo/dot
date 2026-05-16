// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

// BaseCertificate dot
type BaseCertificate struct {
	logger *dot.LoggerType
}

func NewBaseCertificate(logger *dot.LoggerType) *BaseCertificate {
	return &BaseCertificate{
		logger: logger,
	}
}

// GenerateRoot Generate root certificate
func (c *BaseCertificate) GenerateRoot(signatureAlgorithm x509.SignatureAlgorithm, dnsName []string, orgName []string) (*x509.Certificate, error) {
	serialNumber, err := c.makeSerialNumber()
	if err != nil {
		return nil, err
	}

	rootCert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:       []string{"cn"},
			Locality:      []string{"scry"},
			Province:      []string{"scry"},
			Organization:  orgName,
			StreetAddress: []string{"scry"},
			PostalCode:    []string{"scry"},
			CommonName:    "scry",
		},
		DNSNames:              dnsName,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		SignatureAlgorithm:    signatureAlgorithm,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	return rootCert, err
}

// GenerateRootFile
// keyFile private key, pemFile from certificate file
func (c *BaseCertificate) GenerateRootFile(pri any, rootCert *x509.Certificate, pub any, keyFile string, pemFile string) error {
	{
		certBytes, err := x509.CreateCertificate(rand.Reader, rootCert, rootCert, pub, pri)
		if err != nil {
			return err
		}

		file := ""
		file, err = wdFileAndMakeDirs(pemFile)
		if err != nil {
			return err
		}
		certOut, err := os.Create(file)
		if err != nil {
			return err
		}
		defer certOut.Close()
		if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
			return err
		}
	}
	return c.KeyFile(pri, keyFile)
}

// KeyFile
// keyFile private file
func (c *BaseCertificate) KeyFile(leafPri any, keyFile string) error {
	privBytes, err := x509.MarshalPKCS8PrivateKey(leafPri)
	if err != nil {
		return err
	}
	file := ""
	file, err = wdFileAndMakeDirs(keyFile)
	if err != nil {
		return err
	}
	keyOut, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	if err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	return nil
}

// GenerateLeafCertificate Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *BaseCertificate) GenerateLeafCertificate(signatureAlgorithm x509.SignatureAlgorithm, dnsName []string, orgName []string) (*x509.Certificate, error) {
	serialNumber, err := c.makeSerialNumber()
	if err != nil {
		return nil, err
	}
	leafCert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:       []string{"cn"},
			Locality:      []string{"scry"},
			Province:      []string{"scry"},
			Organization:  orgName,
			StreetAddress: []string{"scry"},
			PostalCode:    []string{"scry"},
			CommonName:    "scry",
		},
		DNSNames:              dnsName,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  false,
		SignatureAlgorithm:    signatureAlgorithm,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	return leafCert, nil
}

// GenerateECDSALeafFile Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *BaseCertificate) GenerateLeafFile(leafPri any, leafCert *x509.Certificate, leafPub any, keyFile string, pemFile string, rootCert *x509.Certificate, rootPri any) error {

	{
		certBytes, err := x509.CreateCertificate(rand.Reader, leafCert, rootCert, leafPub, rootPri)
		if err != nil {
			return err
		}

		file := ""
		file, err = wdFileAndMakeDirs(pemFile)
		if err != nil {
			return err
		}
		certOut, err := os.Create(file)
		if err != nil {
			return err
		}
		defer certOut.Close()
		if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
			return err
		}
	}

	return c.KeyFile(leafPri, keyFile)
}

func (c *BaseCertificate) LoadPrivateKey(keyFile string) (any, error) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("is not valide pem file")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		if rsaKey, err2 := x509.ParsePKCS1PrivateKey(block.Bytes); err2 == nil {
			return rsaKey, nil
		}
		return nil, err
	}

	return key, nil
}

func (c *BaseCertificate) KeyType(key any) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		_ = k
	case *ecdsa.PrivateKey:

	case ed25519.PrivateKey:

	case ecdh.PrivateKey:

		// case *mlkem768.PrivateKey: //go 1.27
	}
}

// LoadCertificate Read certificate from pemFile
func (c *BaseCertificate) LoadCertificate(pemFile string) (cert *x509.Certificate, err error) {
	file, err := wdFile(pemFile)

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

	if block.Type == "CERTIFICATE" {
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("must be CERTIFICATE")
	}
	return cert, err
}

func (c *BaseCertificate) ServerName(cert *x509.Certificate) string {
	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0]
	}
	if len(cert.IPAddresses) > 0 {
		return cert.IPAddresses[0].String()
	}
	if len(cert.Subject.CommonName) > 0 {
		return cert.Subject.CommonName
	}
	return ""
}

func (c *BaseCertificate) makeSerialNumber() (serial *big.Int, err error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err = rand.Int(rand.Reader, serialNumberLimit)
	return
}

// If not absolute path, then comparing wd file path, create content where file existing
func wdFileAndMakeDirs(file string) (nfile string, err error) {
	nfile = ""
	tfile := ""
	if filepath.IsAbs(file) {
		tfile = file
	} else {
		exPath, err3 := os.Getwd()
		if err3 != nil {
			err = err3
			return
		}
		tfile = filepath.Join(exPath, file)
	}

	tdir := filepath.Dir(tfile)
	if !sfile.ExistFile(tdir) {
		err = os.MkdirAll(tdir, 0755)
	}
	nfile = tfile
	return
}

// If comparing work path file do not exist, then check whether parameter file existing or not
func wdFile(file string) (tfile string, err error) {
	tfile = ""
	if filepath.IsAbs(file) {
		tfile = file
	} else {
		exPath, err3 := os.Getwd()
		if err3 != nil {
			err = err3
			return
		}
		tfile = filepath.Join(exPath, file)
	}

	if !sfile.ExistFile(tfile) {
		if sfile.ExistFile(file) {
			tfile = file
		} else {
			tfile = ""
			err = errors.New("no file")
		}
	}
	return
}
