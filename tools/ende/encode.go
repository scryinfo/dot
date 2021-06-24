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

const EncodeTypeID = "4f0d40c6-9822-4346-ae23-3a54b866b96a"

type configEncode struct {
	Help     bool   `json:"help"` //输出帮助信息
	Generate string `json:"generate"`
	File     string `json:"file"`
	OutFile  string `json:"outFile"`

	EndeType          string `json:"endeType"`          //加密的实现，现支持X25519
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

	if c.conf.EndeType != string(scrypto.EndeType_X25519) {
		logger.Errorln("the EndeType is not support: " + c.conf.EndeType)
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
		EndeType:        scrypto.EndeType(c.conf.EndeType),
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
		ecdh := scrypto.GetEcdh(&data)
		encoder := scrypto.GetAsymmetricEncoder(&data)
		if encoder == nil {
			logger.Errorln("the EndeType is not support: " + string(data.EndeType))
			return false
		}
		privateKey, _, err := ecdh.GenerateKey(nil)
		if err != nil {
			logger.Errorln("", zap.Error(err))
			return false
		}
		var peerKey crypto.PublicKey
		{
			bytes, err := hex.DecodeString(c.conf.X25519PeerKey)
			if err != nil {
				logger.Errorln("", zap.Error(err))
				return false
			}
			peerKey, err = ecdh.BytesToPublicKey(bytes)
			if err != nil {
				logger.Errorln("", zap.Error(err))
				return false
			}
		}
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
