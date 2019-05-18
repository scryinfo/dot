package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
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

type configEcdsa struct {
	Name string `json:"name"`
}

type Ecdsa struct {
}

func newEcdsa(conf interface{}) (dot.Dot, error) {
	var err error = nil
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

//生成ca证书与私钥
// keyFile 私钥文件， pemFile ca证书文件
func (c *Ecdsa) GenerateCaCertKey(caPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (ca *x509.Certificate, err error) {

	var serialNumber *big.Int = nil
	serialNumber, err = c.makeSerialNumber()

	ca = &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: orgName,
		},
		DNSNames:              dnsName,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(20, 0, 0),
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

//生成子证书与私钥
// keyFile 私钥文件， pemFile 子证书文件
func (c *Ecdsa) GenerateCertKey(caParent *x509.Certificate, caPri *ecdsa.PrivateKey, keyFile string, pemFile string, dnsName []string, orgName []string) (err error) {
	var serialNumber *big.Int = nil
	serialNumber, err = c.makeSerialNumber()
	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: orgName,
		},
		DNSNames:    dnsName,
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(20, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
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

//从keyFile中读取私钥
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

//从pemFile中读取公钥
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

//从pemFile中读取证书
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

//如果不是绝对路径，那么就相对于可执行文件所在路径， 且创建文件所在的目录
func exPathFileAndMakeDirs(file string) (nfile string, err error) {
	nfile = ""
	tfile := ""
	if filepath.IsAbs(file) {
		tfile = file
	} else {
		exPath := ""
		{
			ex, err3 := os.Executable()
			if err == nil {
				exPath = filepath.Dir(ex)
			} else {
				err = err3
				return
			}
		}
		tfile = filepath.Join(exPath, file)
	}

	tdir := filepath.Dir(tfile)
	if !sfile.ExitFile(tdir) {
		err = os.MkdirAll(tdir, os.ModeDir)
	}
	nfile = tfile
	return
}

//如果相对于可执行文件不存在， 就直接看参数file是否存在
func exPathFile(file string) (tfile string, err error) {
	tfile = ""
	if filepath.IsAbs(file) {
		tfile = file
	} else {
		exPath := ""
		{
			ex, err3 := os.Executable()
			if err == nil {
				exPath = filepath.Dir(ex)
			} else {
				err = err3
				return
			}
		}
		tfile = filepath.Join(exPath, file)
	}

	if !sfile.ExitFile(tfile) {
		if sfile.ExitFile(file) {
			tfile = file
		} else {
			tfile = ""
			err = errors.New("no file")
		}
	}
	return
}

//生成私钥
func MakePriKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
