package ende

import (
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/scrypto"
	_ "github.com/scryinfo/dot/lib/scrypto/ende25519"
	_ "github.com/scryinfo/dot/lib/scrypto/sx25519"
	"github.com/scryinfo/dot/utils"
	"github.com/scryinfo/scryg/sutils/sfile"
)

const DecodeTypeID = "7e592fd2-6cd8-4953-8fae-8c3a9915bef7"

type DecodeConfig struct {
	Help     bool   `json:"help"` //输出帮助信息
	Generate string `json:"generate"`
	File     string `json:"file"`
	OutFile  string `json:"outFile"`

	//EndeType string `json:"endeType"` //加密的实现，现支持X25519
	Ed25519PublicKey    string `json:"ed25519PublicKey"`    //用于签名的公钥, hex
	X25519PrivateKeyKey string `json:"x25519PrivateKeyKey"` //解密用的x25519的私钥， hex
}
type Decode struct {
	conf DecodeConfig
}

//func (c *Decode) Create(l dot.Line) error {
//
//}
//func (c *Decode) Injected(l dot.Line) error {
//
//}
//func (c *Decode) AfterAllInject(l dot.Line) {
//
//}

func (c *Decode) Start(ignore bool) error {
	c.parseEnParameter()
	SendSignal(syscall.SIGKILL)
	//os.Exit(0)
	return nil
}

//func (c *Decode) Stop(ignore bool) error {
//
//}
//
//func (c *Decode) Destroy(ignore bool) error {
//
//}

// construct dot
func NewDecode(conf *DecodeConfig) (*Decode, error) {
	d := &Decode{conf: *conf}
	return d, nil
}

// 返回值为true 接着运行
func (c *Decode) parseEnParameter() bool {
	logger := &dot.Logger
	if c.conf.Help {
		//todo
		return false
	}

	if len(c.conf.Generate) > 0 {
		GenerateKey(c.conf.Generate)
		return false
	}

	if len(c.conf.File) < 1 {
		logger.Error().Msg("请入要解密的文件")
		return false
	}

	fullPath := utils.GetFullPathFile(c.conf.File)
	if len(fullPath) < 1 || !sfile.ExistFile(fullPath) {
		logger.Error().Msgf("不能找到文件： %s\n", c.conf.File)
		return false
	}

	if len(c.conf.OutFile) < 1 {
		c.conf.OutFile = fullPath + "_de"
	}

	data := scrypto.EndeData{}
	{
		bytes, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error().Err(err).Send()
			return false
		}
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			logger.Error().Err(err).Send()
			return false
		}
	}
	{ //decode
		var privateKey crypto.PrivateKey
		{
			bytes, err := hex.DecodeString(c.conf.X25519PrivateKeyKey)
			if err != nil {
				logger.Error().Err(err).Send()
				return false
			}
			privateKey = bytes
		}
		var signedPublicKey ed25519.PublicKey
		{
			bytes, err := hex.DecodeString(c.conf.Ed25519PublicKey)
			if err != nil {
				logger.Error().Err(err).Send()
				return false
			}
			signedPublicKey = bytes
		}
		var err error
		data, err = scrypto.DecodeData(&data, privateKey, signedPublicKey)
		if err != nil {
			logger.Error().Err(err).Send()
			return false
		}
	}

	err := os.WriteFile(c.conf.OutFile, data.Body, 0644)
	if err != nil {
		logger.Error().Err(err).Send()
		return false
	}

	fmt.Printf("decode successful, file: %s\n", c.conf.OutFile)

	return true
}
