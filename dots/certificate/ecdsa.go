// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

const (
	EcdsaTypeId = "4b8b1751-4799-4578-af46-d9b339cf582f"
)

type Ecdsa struct {
}

func newEcdsa(conf interface{}) (dot.Dot, error) {
	var err error = nil
	_ = conf
	d := &Ecdsa{}
	return d, err
}

func TypeLiveEcdsa() *dot.TypeLives {
	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: EcdsaTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
			return newEcdsa(conf)
		}},
	}
}

//Generate ca certificate and private key
// keyFile private key, pemFile ca certificate file
func (c *Ecdsa) GenerateCaCertKey(caPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (ca *x509.Certificate, err error) {

	var serialNumber *big.Int = nil
	serialNumber, err = c.makeSerialNumber()

	ca = &x509.Certificate{
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
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	{
		var certBytes []byte = nil
		pub := &caPri.PublicKey
		certBytes, err = x509.CreateCertificate(rand.Reader, ca, ca, pub, caPri)
		if err != nil {
			return nil, err
		}

		file := ""
		file, err = exPathFileAndMakeDirs(pemFile)
		if err != nil {
			return nil, err
		}
		certOut, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		defer certOut.Close()
		if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
			return nil, err
		}
	}
	{
		var privBytes []byte = nil
		if privBytes, err = x509.MarshalECPrivateKey(caPri); err != nil {
			return nil, err
		}
		file := ""
		file, err = exPathFileAndMakeDirs(keyFile)
		if err != nil {
			return nil, err
		}
		keyOut, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, err
		}
		defer keyOut.Close()
		if err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
			return nil, err
		}
	}

	return ca, err
}

//Generate subcertificate and private key
// keyFile private file, pemFile subcertificate file
func (c *Ecdsa) GenerateCertKey(caParent *x509.Certificate, caPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (err error) {
	var serialNumber *big.Int = nil
	serialNumber, err = c.makeSerialNumber()
	cert := &x509.Certificate{
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
		DNSNames:           dnsName,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().AddDate(100, 0, 0),
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	{
		pub := &priv.PublicKey
		certBytes, err := x509.CreateCertificate(rand.Reader, cert, caParent, pub, caPri)
		if err != nil {
			return err
		}

		file := ""
		file, err = exPathFileAndMakeDirs(pemFile)
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

	{
		var privBytes []byte = nil
		if privBytes, err = x509.MarshalECPrivateKey(priv); err != nil {
			return err
		}
		file := ""
		file, err = exPathFileAndMakeDirs(keyFile)
		if err != nil {
			return err
		}
		keyOut, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer keyOut.Close()
		if err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
			return err
		}
	}

	return err
}

//Read private key from keyFile
func (c *Ecdsa) PrivateKey(keyFile string) (pri *ecdsa.PrivateKey, err error) {

	file, err := exPathFile(keyFile)

	if err != nil || len(file) < 1 {
		return nil, errors.New("file do not exist")
	}

	bs, err := ioutil.ReadFile(file)
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

//Read public key from pemFile
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

//Read certificate from pemFile
func (c *Ecdsa) Certificate(pemFile string) (cert *x509.Certificate, err error) {
	file, err := exPathFile(pemFile)

	if err != nil || len(file) < 1 {
		return nil, errors.New("file do not exist")
	}

	bs, err := ioutil.ReadFile(file)
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

func (c *Ecdsa) makeSerialNumber() (serial *big.Int, err error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err = rand.Int(rand.Reader, serialNumberLimit)
	return
}

//If not absolute path, then comparing executable file path, create content where file existing
func exPathFileAndMakeDirs(file string) (nfile string, err error) {
	nfile = ""
	tfile := ""
	if filepath.IsAbs(file) {
		tfile = file
	} else {
		exPath := ""
		{
			ex, err3 := os.Executable()
			if err3 == nil {
				exPath = filepath.Dir(ex)
			} else {
				err = err3
				return
			}
		}
		tfile = filepath.Join(exPath, file)
	}

	tdir := filepath.Dir(tfile)
	if !sfile.ExistFile(tdir) {
		err = os.MkdirAll(tdir, os.ModeDir)
	}
	nfile = tfile
	return
}

//If comparing executable file do not exist, then check whether parameter file existing or not
func exPathFile(file string) (tfile string, err error) {
	tfile = ""
	if filepath.IsAbs(file) {
		tfile = file
	} else {
		exPath := ""
		{
			ex, err3 := os.Executable()
			if err3 == nil {
				exPath = filepath.Dir(ex)
			} else {
				err = err3
				return
			}
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

//Generate private key
func MakePriKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
