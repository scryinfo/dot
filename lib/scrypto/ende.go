package scrypto

import (
	"crypto"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
)

type AsymmetricDecoder interface {
	// EcdhDecode
	// EndeData中已经包含peersKey，所以没有出现在参数中
	EcdhDecode(privateKey crypto.PrivateKey, cipher *EndeData) (plain EndeData, err error)
}

type AsymmetricEncoder interface {
	EcdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plain *EndeData) (cipher EndeData, err error)
}

type EndeType string

const (
	EndeType_X25519 EndeType = "x25519"
)

// 以下值只会在init中进行修改，所以认为是安全的
var (
	Encoders = make(map[EndeType]AsymmetricEncoder, 2)
	Decoders = make(map[EndeType]AsymmetricDecoder, 2)
)

type EndeData struct {
	//加密时使用的公钥
	PublicKey []byte
	EndeType  EndeType
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

func toJsonEndeData(data *EndeData) (to endeData) {
	to.PublicKey = hex.EncodeToString(data.PublicKey)
	to.EndeType = string(data.EndeType)
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
	c.EndeType = EndeType(data.EndeType)
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

func GetAsymmetricDecoder(data *EndeData) AsymmetricDecoder {
	var re AsymmetricDecoder = nil
	if v, ok := Decoders[data.EndeType]; ok {
		re = v
	}
	return re
}

func GetAsymmetricEncoder(data *EndeData) AsymmetricEncoder {
	var re AsymmetricEncoder = nil
	if v, ok := Encoders[data.EndeType]; ok {
		re = v
	}
	return re
}

func ToSigningHash(data []byte) (hash []byte) {
	h251 := sha256.New()
	if len(data) > 0 {
		h251.Write(data)
	}
	first := h251.Sum(nil)
	h251.Reset()
	h251.Write(first) //两次hash
	return h251.Sum(nil)
}

func SignEd25519(privateKey ed25519.PrivateKey, hash []byte) (signature []byte, publicBytes []byte) {
	signature = ed25519.Sign(privateKey, hash)
	publicBytes = privateKey.Public().(ed25519.PublicKey)
	return
}

func EncodeData(plain *EndeData, peerKey crypto.PublicKey, signedKey ed25519.PrivateKey) (cipher EndeData, err error) {
	cipher = *plain
	if signedKey != nil {
		cipher.Signature, cipher.SignedPublicKey = SignEd25519(signedKey, ToSigningHash(plain.Body))
	}
	ecdh := GetEcdh(&cipher)
	privateKey, _, err := ecdh.GenerateKey(nil)
	if err != nil {
		return
	}
	encoder := GetAsymmetricEncoder(&cipher)
	cipher, err = encoder.EcdhEncode(privateKey, peerKey, &cipher)
	if err != nil {
		return
	}

	return
}
func DecodeData(cipher *EndeData, privateKey crypto.PrivateKey, signedPublicKey ed25519.PublicKey) (plain EndeData, err error) {
	plain = *cipher
	encoder := GetAsymmetricDecoder(&plain)
	plain, err = encoder.EcdhDecode(privateKey, &plain)
	if err != nil {
		return
	}

	if signedPublicKey != nil { //verify
		hash := ToSigningHash(plain.Body)
		if !signedPublicKey.Equal(ed25519.PublicKey(plain.SignedPublicKey)) || !ed25519.Verify(plain.SignedPublicKey, hash, plain.Signature) {
			err = errors.New("failed to verify signature")
			return
		}
	}

	return
}
