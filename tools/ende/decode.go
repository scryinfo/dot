package ende

import (
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"syscall"

	"go.uber.org/zap"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/scrypto"
	_ "github.com/scryinfo/dot/lib/scrypto/ende25519"
	_ "github.com/scryinfo/dot/lib/scrypto/sx25519"
	"github.com/scryinfo/dot/utils"
	"github.com/scryinfo/scryg/sutils/sfile"
)

const DecodeTypeID = "7e592fd2-6cd8-4953-8fae-8c3a9915bef7"

type configDecode struct {
	Help     bool   `json:"help"` //输出帮助信息
	Generate string `json:"generate"`
	File     string `json:"file"`
	OutFile  string `json:"outFile"`

	//EndeType string `json:"endeType"` //加密的实现，现支持X25519
	Ed25519PublicKey    string `json:"ed25519PublicKey"`    //用于签名的公钥, hex
	X25519PrivateKeyKey string `json:"x25519PrivateKeyKey"` //解密用的x25519的私钥， hex
}
type Decode struct {
	conf configDecode
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

//construct dot
func newDecode(conf []byte) (dot.Dot, error) {
	dconf := &configDecode{}

	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &Decode{conf: *dconf}

	return d, nil
}

//DecodeTypeLives
func DecodeTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: DecodeTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newDecode(conf)
		}},
		//Lives: []dot.Live{
		//	{
		//		LiveID:    DecodeTypeID,
		//		RelyLives: map[string]dot.LiveID{"some field": "some id"},
		//	},
		//},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//DecodeConfigTypeLive
func DecodeConfigTypeLive() *dot.ConfigTypeLive {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLive{
		TypeIDConfig: DecodeTypeID,
		ConfigInfo: &configDecode{
			//todo
		},
	}
}

//返回值为true 接着运行
func (c *Decode) parseEnParameter() bool {
	logger := dot.Logger()
	if c.conf.Help {
		//todo
		return false
	}

	if len(c.conf.Generate) > 0 {
		GenerateKey(c.conf.Generate)
		return false
	}

	if len(c.conf.File) < 1 {
		logger.Errorln("请入要解密的文件")
		return false
	}

	fullPath := utils.GetFullPathFile(c.conf.File)
	if len(fullPath) < 1 || !sfile.ExistFile(fullPath) {
		logger.Errorln(fmt.Sprintf("不能找到文件： %s\n", c.conf.File))
		return false
	}

	if len(c.conf.OutFile) < 1 {
		c.conf.OutFile = fullPath + "_de"
	}

	data := scrypto.EndeData{}
	{
		bytes, err := ioutil.ReadFile(fullPath)
		if err != nil {
			logger.Errorln("", zap.Error(err))
			return false
		}
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			logger.Errorln("", zap.Error(err))
			return false
		}
	}
	{ //decode
		var privateKey crypto.PrivateKey
		{
			bytes, err := hex.DecodeString(c.conf.X25519PrivateKeyKey)
			if err != nil {
				logger.Errorln("", zap.Error(err))
				return false
			}
			privateKey = bytes
		}
		var signedPublicKey ed25519.PublicKey
		{
			bytes, err := hex.DecodeString(c.conf.Ed25519PublicKey)
			if err != nil {
				logger.Errorln("", zap.Error(err))
				return false
			}
			signedPublicKey = bytes
		}
		var err error
		data, err = scrypto.DecodeData(&data, privateKey, signedPublicKey)
		if err != nil {
			logger.Errorln("", zap.Error(err))
			return false
		}
	}

	err := ioutil.WriteFile(c.conf.OutFile, data.Body, 0644)
	if err != nil {
		logger.Errorln("", zap.Error(err))
		return false
	}

	fmt.Printf("decode successful, file: %s\n", c.conf.OutFile)

	return true
}
