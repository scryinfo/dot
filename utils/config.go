package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

func GetFullPathFile(file string) string {
	if filepath.IsAbs(file) {
		return file
	}

	res := ""
	for {
		ex, err := os.Executable()
		if err != nil {
			dot.Logger().Errorln("connsImp", zap.Error(err))
			res = ""
			break
		}

		ex = filepath.Dir(ex)
		temp := filepath.Join(ex, file)
		if sfile.ExistFile(temp) {
			res = temp
			break
		} else { //try find file from the current path
			temp, err = os.Getwd()
			if err != nil {
				dot.Logger().Errorln("connsImp", zap.Error(err))
				res = ""
				break
			}
			temp = filepath.Join(temp, file)
			if sfile.ExistFile(temp) {
				res = temp
				break
			}
		}

		break
	}

	return res

}
func GetTlsConfig(conf *TlsConfig) (*tls.Config, error) {
	if len(conf.CaPem) > 0 && len(conf.Key) > 0 && len(conf.Pem) > 0 { //both tls
		caPemFile := GetFullPathFile(conf.CaPem)
		if len(caPemFile) < 1 {
			return nil, errors.New("the caPem is not empty, and can not find the file: " + conf.CaPem)
		}
		keyFile := GetFullPathFile(conf.Key)
		if len(keyFile) < 1 {
			return nil, errors.New("the Key is not empty, and can not find the file: " + conf.Key)
		}

		pemFile := GetFullPathFile(conf.Pem)
		if len(pemFile) < 1 {
			return nil, errors.New("the Pem is not empty, and can not find the file: " + conf.Pem)
		}

		pool := x509.NewCertPool()
		{
			caCrt, err := ioutil.ReadFile(caPemFile)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if !pool.AppendCertsFromPEM(caCrt) {
				return nil, errors.New("credentials: failed to append certificates")
			}
		}
		cert, err := tls.LoadX509KeyPair(pemFile, keyFile)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return &tls.Config{
			ServerName:   conf.ServerNameOverride,
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool, //ClientConn, use the RootCAs
		}, nil

	} else if len(conf.Pem) > 0 { //just server
		pemFile := GetFullPathFile(conf.Pem)
		if len(pemFile) < 1 {
			return nil, errors.New("the Pem is not empty, and can not find the file: " + conf.Pem)
		}

		b, err := ioutil.ReadFile(pemFile)
		if err != nil {
			return nil, err
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}

		return &tls.Config{ServerName: conf.ServerNameOverride, RootCAs: cp}, nil
	}
	return nil, nil
}
