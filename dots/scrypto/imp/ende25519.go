package imp

import (
	"crypto"
	"crypto/sha256"
	"github.com/scryinfo/dot/dots/scrypto"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"io"
)

type ende25519 struct{}

func EcdhDecoder25519() scrypto.EcdhDecoder {
	return &ende25519{}
}

func EcdhEncoder25519() scrypto.EcdhEncoder {
	return &ende25519{}
}

var (
	salt  []byte
	hash  = sha256.New
	info  = []byte("scry info")
	nonce = make([]byte, chacha20poly1305.NonceSize)
	ecdh  = X25519()
)

func (c *ende25519) EcdhDecode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, ciphertext []byte) (plaintext []byte, err error) {

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

func (c *ende25519) EcdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plaintext []byte) (ciphertext []byte, err error) {
	echg := X25519()
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
