package ende

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sx255192 "github.com/scryinfo/dot/lib/scrypto/sx25519"
	"io"
	"log"
)

const (
	generate_ed25519 = "ed25519"
	generate_x25519  = "x25519"
)

func GenerateEd25519Key(rand io.Reader) (privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, err error) {
	publicKey, privateKey, err = ed25519.GenerateKey(rand)
	return
}

func GenerateX25519Key(rand io.Reader) (privateKey sx255192.PrivateKey, publicKey sx255192.PublicKey, err error) {
	x := sx255192.X25519()
	privateKey2, publicKey2, err := x.GenerateKey(rand)
	if err != nil {
		return
	}

	err = privateKey.ToPrivateKey(privateKey2.([]byte))
	if err != nil {
		return
	}

	err = publicKey.ToPublicKey(publicKey2.([]byte))
	//if err != nil {
	//	return
	//}

	return
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

func VerifyEd25519(publicKey ed25519.PublicKey, hash []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, hash, signature)
}

func GenerateKey(name string) {
	switch name {
	case generate_ed25519:
		privateKey, publicKey, err := GenerateEd25519Key(nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ed25519 keys\nprivate: %s\npublic: %s\n", hex.EncodeToString(privateKey), hex.EncodeToString(publicKey))
		break
	case generate_x25519:
		privateKey, publicKey, err := GenerateX25519Key(nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("x25519 keys\nprivate: %s\npublic: %s\n", hex.EncodeToString(privateKey), hex.EncodeToString(publicKey))
		break
	default:
		break
	}
}
