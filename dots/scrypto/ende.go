package scrypto

import (
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

type EcdhDecoder interface {
	// EcdhDecode
	// EndeData中已经包含peersKey，所以没有出现在参数中
	EcdhDecode(privateKey crypto.PrivateKey, cipher EndeData) (plain EndeData, err error)
}

type EcdhEncoder interface {
	EcdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plain EndeData) (cipher EndeData, err error)
}

const (
	EndeType_X25519 = "X25519"
)

type EndeData struct {
	//加密时使用的公钥
	PublicKey []byte
	EndeType  string
	//对Body的签名，
	Signature []byte
	//ed25519不能从签名出计算出公钥，所以需要来验证签名。在使用中双方也可以约定公钥，并不一定要存放在字段中
	SignedPublicKey []byte
	//true数据已加密
	EnData bool
	//数据
	Body []byte
}

type endeData struct {
	PublicKey       string //hex
	EndeType        string
	Signature       string //hex
	SignedPublicKey string //hex
	EnData          bool
	Body            string //base64
}

func (c *EndeData) UnmarshalJSON(b []byte) error {
	var data endeData
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	return c.toEndeData(&data)
}

func (c EndeData) MarshalJSON() ([]byte, error) {
	to := toJsonEndeData(&c)
	return json.Marshal(to)
}

//生成hash，并记录下来
func (c *EndeData) toSigningHash() []byte {
	h251 := sha256.New()
	if len(c.Body) > 0 {
		h251.Write(c.Body)
	}
	first := h251.Sum(nil)
	h251.Reset()
	h251.Write(first) //两次hash
	return h251.Sum(nil)
}

func toJsonEndeData(data *EndeData) (to endeData) {
	to.PublicKey = hex.EncodeToString(data.PublicKey)
	to.EndeType = data.EndeType
	to.Signature = hex.EncodeToString(data.Signature)
	to.SignedPublicKey = hex.EncodeToString(data.SignedPublicKey)
	to.EnData = data.EnData
	to.Body = base64.StdEncoding.EncodeToString(data.Body)
	return
}

func (c *EndeData) toEndeData(data *endeData) (err error) {
	c.PublicKey, err = hex.DecodeString(data.PublicKey)
	if err != nil {
		return
	}
	c.EndeType = data.EndeType
	if err != nil {
		return
	}
	c.Signature, err = hex.DecodeString(data.Signature)
	if err != nil {
		return
	}
	c.SignedPublicKey, err = hex.DecodeString(data.SignedPublicKey)
	if err != nil {
		return
	}
	c.EnData = data.EnData
	c.Body, err = base64.StdEncoding.DecodeString(data.Body)
	if err != nil {
		return
	}
	return
}
