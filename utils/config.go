package utils

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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
