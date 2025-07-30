package ende25519

import (
	"crypto"
	"crypto/sha256"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/lib/scrypto"
	"github.com/scryinfo/dot/lib/scrypto/sx25519"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"io"
)

type ende25519 struct{}

func EcdhDecoder25519() scrypto.AsymmetricDecoder {
	return &ende25519{}
}

func EcdhEncoder25519() scrypto.AsymmetricEncoder {
	return &ende25519{}
}

var (
	salt     []byte
	hash     = sha256.New
	info     = []byte("scry info")
	nonce    = make([]byte, chacha20poly1305.NonceSize)
	ecdh     = sx25519.X25519()
	endeType = scrypto.EndeType_X25519
)

func init() {
	scrypto.Encoders[scrypto.EndeType_X25519] = EcdhEncoder25519()
	scrypto.Decoders[scrypto.EndeType_X25519] = EcdhDecoder25519()
}

func (c *ende25519) EcdhDecode(privateKey crypto.PrivateKey, _cipher *scrypto.EndeData) (plain scrypto.EndeData, err error) {
	plain = *_cipher
	if !plain.EnData {
		return
	}
	if plain.EndeType != endeType {
		err = errors.New("the ende type is not " + string(endeType))
		return
	}

	var peersKey crypto.PublicKey = plain.PublicKey

	plain.Body, err = c._ecdhDecode(privateKey, peersKey, plain.Body)
	if err != nil {
		return
	}
	plain.EnData = false //decode data

	return
}

func (c *ende25519) EcdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, _plain *scrypto.EndeData) (cipher scrypto.EndeData, err error) {
	cipher = *_plain
	if cipher.EnData {
		return cipher, nil
	}
	cipher.EndeType = endeType
	publicKey, err := ecdh.PublicKey(privateKey)
	if err != nil {
		return
	}
	cipher.PublicKey, err = ecdh.PublicKeyToBytes(publicKey)
	if err != nil {
		return
	}

	cipher.Body, err = c._ecdhEncode(privateKey, peersKey, cipher.Body)
	if err != nil {
		return
	}
	cipher.EnData = true
	return
}

// EcdhDecode
// privateKey sx25519, peersKey sx25519
func (c *ende25519) _ecdhDecode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, ciphertext []byte) (plaintext []byte, err error) {

	key, err := ecdh.ComputeSecret(privateKey, peersKey)
	if err != nil {
		return
	}
	dk := hkdf.New(hash, key, salt, info)
	wrappingKey := make([]byte, chacha20poly1305.KeySize)
	if _, err = io.ReadFull(dk, wrappingKey); err != nil {
		return
	}
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}
	plaintext, err = aead.Open(nil, nonce, ciphertext, nil)
	return
}

// EcdhEncode
// privateKey sx25519, peersKey sx25519
func (c *ende25519) _ecdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plaintext []byte) (ciphertext []byte, err error) {
	echg := sx25519.X25519()
	key, err := echg.ComputeSecret(privateKey, peersKey)
	if err != nil {
		return
	}

	dk := hkdf.New(hash, key, salt, info)

	wrappingKey := make([]byte, chacha20poly1305.KeySize)
	if _, err = io.ReadFull(dk, wrappingKey); err != nil {
		return
	}

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}
	ciphertext = aead.Seal(nil, nonce, plaintext, nil)
	return
}
