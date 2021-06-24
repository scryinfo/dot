package ende

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/scrypto"
	"github.com/scryinfo/dot/dots/scrypto/ende25519"
	"github.com/scryinfo/dot/dots/scrypto/sx25519"
	"github.com/scryinfo/dot/utils"
	"github.com/scryinfo/scryg/sutils/sfile"
	"go.uber.org/zap"
	"io/ioutil"
	"syscall"
)

const EncodeTypeID = "4f0d40c6-9822-4346-ae23-3a54b866b96a"

type configEncode struct {
	Help     bool   `json:"help"` //输出帮助信息
	Generate string `json:"generate"`
	File     string `json:"file"`
	OutFile  string `json:"outFile"`

	Ed25519PrivateKey string `json:"ed25519PrivateKey"` //用于签名的key, hex
	X25519PeerKey     string `json:"x25519PeerKey"`     //对方的x25519的公钥， hex
}
type Encode struct {
	conf configEncode
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

//construct dot
func newEncode(conf []byte) (dot.Dot, error) {
	dconf := &configEncode{}

	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &Encode{conf: *dconf}

	return d, nil
}

//EncodeTypeLives
func EncodeTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: EncodeTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newEncode(conf)
		}},
		//Lives: []dot.Live{
		//	{
		//		LiveID:    EncodeTypeID,
		//		RelyLives: map[string]dot.LiveID{"some field": "some id"},
		//	},
		//},
	}

	lives := []*dot.TypeLives{tl}

	return lives
}

//EncodeConfigTypeLive
func EncodeConfigTypeLive() *dot.ConfigTypeLive {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLive{
		TypeIDConfig: EncodeTypeID,
		ConfigInfo: &configEncode{
			//todo
		},
	}
}

//返回值为true 接着运行
func (c *Encode) parseEnParameter() bool {
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
		logger.Errorln("请入要加密的文件")
		return false
	}

	fullPath := utils.GetFullPathFile(c.conf.File)
	if len(fullPath) < 1 || !sfile.ExistFile(fullPath) {
		logger.Errorln(fmt.Sprintf("不能找到文件： %s\n", c.conf.File))
		return false
	}

	if len(c.conf.OutFile) < 1 {
		c.conf.OutFile = fullPath + "_en"
	}

	body, err := ioutil.ReadFile(fullPath)
	if err != nil {
		logger.Errorln("", zap.Error(err))
		return false
	}

	data := scrypto.EndeData{
		PublicKey:       nil,
		EndeType:        "",
		Signature:       nil,
		SignedPublicKey: nil,
		EnData:          false,
		Body:            body,
	}
	{
		var skey ed25519.PrivateKey
		{
			bytes, err := hex.DecodeString(c.conf.Ed25519PrivateKey)
			if err != nil {
				logger.Errorln("", zap.Error(err))
				return false
			}
			skey = bytes
		}
		hash := ToSigningHash(data.Body)
		data.Signature, data.SignedPublicKey = SignEd25519(skey, hash)
	}
	{ //encode
		privateKey, _, err := GenerateX25519Key(nil)
		if err != nil {
			logger.Errorln("", zap.Error(err))
			return false
		}
		var peerKey sx25519.PublicKey
		{
			bytes, err := hex.DecodeString(c.conf.X25519PeerKey)
			if err != nil {
				logger.Errorln("", zap.Error(err))
				return false
			}
			peerKey = bytes
		}
		encoder := ende25519.EcdhEncoder25519()
		data, err = encoder.EcdhEncode(privateKey, peerKey, data)
		if err != nil {
			logger.Errorln("", zap.Error(err))
			return false
		}
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		logger.Errorln("", zap.Error(err))
		return false
	}

	err = ioutil.WriteFile(c.conf.OutFile, bytes, 0644)
	if err != nil {
		logger.Errorln("", zap.Error(err))
		return false
	}
	fmt.Printf("encode successful, file: %s\n", c.conf.OutFile)
	return true
}
