// Scry Info.  All rights reserved.
// license that can be found in the license file.

package shared

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

//TlsConfig
//case the CaPem Pem Key are not empty, both tls
//case only the Pem is not empty, server tls,  if the ServerNameOverride exist, then set the ServerNameOverride that it make the pam and key
type TlsConfig struct {
	//public key of ca
	CaPem string `json:"caPem"`
	//public key of client or server
	Pem string `json:"pem"`
	//private key of client
	Key string `json:"key"`
	//if the CaPam and Key are empty, set the ServerNameOverride is the name of the pem
	ServerNameOverride string `json:"serverNameOverride"`
}

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
