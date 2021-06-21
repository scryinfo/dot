package scrypto

import "crypto"

type EcdhDecoder interface {
	EcdhDecode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, ciphertext []byte) (plaintext []byte, err error)
}

type EcdhEncoder interface {
	EcdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plaintext []byte) (ciphertext []byte, err error)
}
