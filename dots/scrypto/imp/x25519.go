package imp

import (
	"crypto"
	"crypto/rand"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dots/scrypto"
	"golang.org/x/crypto/curve25519"
	"io"
)

//see: https://github.com/aead/ecdh
type ecdh25519 struct{}

var curve25519Params = scrypto.CurveParameters{
	Name:    "Curve25519",
	BitSize: 255,
}

func X25519() scrypto.Ecdh {
	return ecdh25519{}
}

func (ecdh25519) GenerateKey(random io.Reader) (privateKey crypto.PrivateKey, publicKey crypto.PublicKey, err error) {
	if random == nil {
		random = rand.Reader
	}

	var pri, pub [32]byte
	_, err = io.ReadFull(random, pri[:])
	if err != nil {
		return
	}

	// From https://cr.yp.to/ecdh.html
	pri[0] &= 248
	pri[31] &= 127
	pri[31] |= 64

	curve25519.ScalarBaseMult(&pub, &pri)

	privateKey = pri
	publicKey = pub
	return
}

func (ecdh25519) Parameters() scrypto.CurveParameters { return curve25519Params }

func (ecdh25519) PublicKey(private crypto.PrivateKey) (publicKey crypto.PublicKey, err error) {
	var pri, pub [32]byte
	if ok := checkType(&pri, private); !ok {
		return nil, errors.New("ecdh: unexpected type of private key")
	}

	curve25519.ScalarBaseMult(&pub, &pri)

	publicKey = pub
	return
}

func (ecdh25519) Check(peersPublic crypto.PublicKey) (err error) {
	if ok := checkType(new([32]byte), peersPublic); !ok {
		err = errors.New("unexpected type of peers public key")
	}
	return
}

func (ecdh25519) ComputeSecret(privateKey crypto.PrivateKey, peersPublic crypto.PublicKey) (secret []byte, err error) {
	var pri, pub [32]byte
	if ok := checkType(&pri, privateKey); !ok {
		err = errors.New("ecdh: unexpected type of privateKey key")
		return
	}
	if ok := checkType(&pub, peersPublic); !ok {
		err = errors.New("ecdh: unexpected type of peers public key")
		return
	}
	var sec []byte
	sec, err = curve25519.X25519(pri[:], pub[:])
	//curve25519.ScalarMult(&sec, &pri, &pub)

	secret = sec[:]
	return
}

func (ecdh25519) PublicKeyToBytes(publicKey crypto.PublicKey) (key []byte, err error) {
	var bytes [32]byte
	if ok := checkType(&bytes, publicKey); !ok {
		err = errors.New("unexpected type of public key")
	}
	return bytes[:], nil
}

func (ecdh25519) PrivateKeyToBytes(privateKey crypto.PrivateKey) (key []byte, err error) {
	var bytes [32]byte
	if ok := checkType(&bytes, privateKey); !ok {
		err = errors.New("unexpected type of private key")
	}
	return bytes[:], nil
}

func checkType(key *[32]byte, typeToCheck interface{}) (ok bool) {
	switch t := typeToCheck.(type) {
	case [32]byte:
		copy(key[:], t[:])
		ok = true
	case *[32]byte:
		copy(key[:], t[:])
		ok = true
	case []byte:
		if len(t) == 32 {
			copy(key[:], t)
			ok = true
		}
	case *[]byte:
		if len(*t) == 32 {
			copy(key[:], *t)
			ok = true
		}
	}
	return
}
