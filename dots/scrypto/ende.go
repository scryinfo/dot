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
	//签名的hash256
	Hash []byte
	//对没有加密的数据签名，
	Signature []byte
	//true数据已加密
	EnData bool
	//数据
	Body []byte
}

type endeData struct {
	PublicKey string //hex
	EndeType  string
	//签名的hash256
	Hash      string //hex
	Signature string //hex
	EnData    bool
	Body      string //base64
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
func (c *EndeData) toHash() []byte {
	h251 := sha256.New()
	if len(c.PublicKey) > 0 {
		h251.Write(c.PublicKey)
	}
	if len(c.EndeType) > 0 {
		h251.Write([]byte(c.EndeType))
	}
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
	to.Hash = hex.EncodeToString(data.Hash)
	to.Signature = hex.EncodeToString(data.Signature)
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
	c.Hash, err = hex.DecodeString(data.Hash)
	if err != nil {
		return
	}
	c.Signature, err = hex.DecodeString(data.Signature)
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
