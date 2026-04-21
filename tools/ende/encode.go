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

const EncodeTypeID = "4f0d40c6-9822-4346-ae23-3a54b866b96a"

type EncodeConfig struct {
	Help     bool   `json:"help"` //输出帮助信息
	Generate string `json:"generate"`
	File     string `json:"file"`
	OutFile  string `json:"outFile"`

	EndeType          string `json:"endeType"`          //加密的实现，现支持X25519
	Ed25519PrivateKey string `json:"ed25519PrivateKey"` //用于签名的key, hex
	X25519PeerKey     string `json:"x25519PeerKey"`     //对方的x25519的公钥， hex
}
type Encode struct {
	conf EncodeConfig
}

//func (c *Encode) Create(l dot.Line) error {
//
//}
//func (c *Encode) Injected(l dot.Line) error {
//
//}
//func (c *Encode) AfterAllInject(l dot.Line) {
//
//}

func (c *Encode) Start(ignore bool) error {
	c.parseEnParameter()
	SendSignal(syscall.SIGKILL)
	//os.Exit(0)
	return nil
}

//func (c *Encode) Stop(ignore bool) error {
//
//}
//
//func (c *Encode) Destroy(ignore bool) error {
//
//}

// construct dot
func NewEncode(conf *EncodeConfig) *Encode {
	d := &Encode{conf: *conf}
	return d
}

// 返回值为true 接着运行
func (c *Encode) parseEnParameter() bool {
	logger := &dot.Logger
	if c.conf.Help {
		//todo
		return false
	}

	if len(c.conf.Generate) > 0 {
		GenerateKey(c.conf.Generate)
		return false
	}

	if c.conf.EndeType != string(scrypto.EndeType_X25519) {
		logger.Error().Msg("the EndeType is not support: " + c.conf.EndeType)
		return false
	}

	if len(c.conf.File) < 1 {
		logger.Error().Msg("请入要加密的文件")
		return false
	}

	fullPath := utils.GetFullPathFile(c.conf.File)
	if len(fullPath) < 1 || !sfile.ExistFile(fullPath) {
		logger.Error().Msgf("不能找到文件： %s\n", c.conf.File)
		return false
	}

	if len(c.conf.OutFile) < 1 {
		c.conf.OutFile = fullPath + "_en"
	}

	body, err := os.ReadFile(fullPath)
	if err != nil {
		logger.Error().Err(err).Send()
		return false
	}

	data := scrypto.EndeData{
		PublicKey:       nil,
		EndeType:        scrypto.EndeType(c.conf.EndeType),
		Signature:       nil,
		SignedPublicKey: nil,
		EnData:          false,
		Body:            body,
	}

	{ //encode
		var peerKey crypto.PublicKey
		{
			bytes, err := hex.DecodeString(c.conf.X25519PeerKey)
			if err != nil {
				logger.Error().Err(err).Send()
				return false
			}
			peerKey = bytes
		}
		var signedPrivateKey ed25519.PrivateKey
		{
			bytes, err := hex.DecodeString(c.conf.Ed25519PrivateKey)
			if err != nil {
				logger.Error().Err(err).Send()
				return false
			}
			if len(bytes) > 0 {
				signedPrivateKey = bytes
			}
		}
		data, err = scrypto.EncodeData(&data, peerKey, signedPrivateKey)
		if err != nil {
			logger.Error().Err(err).Send()
			return false
		}
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		logger.Error().Err(err).Send()
		return false
	}

	err = os.WriteFile(c.conf.OutFile, bytes, 0644)
	if err != nil {
		logger.Error().Err(err).Send()
		return false
	}
	fmt.Printf("encode successful, file: %s\n", c.conf.OutFile)
	return true
}
