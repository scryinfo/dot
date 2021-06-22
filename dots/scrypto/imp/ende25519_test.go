package imp

import (
	"github.com/scryinfo/dot/dots/scrypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnde25519(t *testing.T) {
	plaintext := []byte("test data")

	ecdh := X25519()
	privateA, publicA, err := ecdh.GenerateKey(nil)
	assert.Equal(t, nil, err)
	privateB, publicB, err := ecdh.GenerateKey(nil)
	assert.Equal(t, nil, err)

	ende := ende25519{}

	ciphertext, err := ende._ecdhEncode(privateA, publicB, plaintext)
	assert.Equal(t, nil, err)

	decodeText, err := ende._ecdhDecode(privateB, publicA, ciphertext)
	assert.Equal(t, nil, err)

	assert.Equal(t, plaintext, decodeText)
}

func TestEnde25519_EcdhDecode(t *testing.T) {

	ecdh := X25519()
	privateA, publicA, err := ecdh.GenerateKey(nil)
	assert.Equal(t, nil, err)
	privateB, publicB, err := ecdh.GenerateKey(nil)
	assert.Equal(t, nil, err)

	ende := ende25519{}

	planData := scrypto.EndeData{
		PublicKey: nil,
		EndeType:  "",
		Hash:      nil,
		Signature: nil,
		EnData:    false,
		Body:      []byte("test struct"),
	}
	cipher, err := ende.EcdhEncode(privateA, publicB, planData)
	assert.Equal(t, nil, err)
	assert.Equal(t, cipher.EndeType, endeType)
	publicABytes, err := ecdh.PublicKeyToBytes(publicA)
	assert.Equal(t, nil, err)
	assert.Equal(t, cipher.PublicKey, publicABytes)
	assert.Equal(t, true, cipher.EnData)

	planData2, err := ende.EcdhDecode(privateB, cipher)
	assert.Equal(t, nil, err)
	assert.Equal(t, planData.Body, planData2.Body)
	assert.Equal(t, planData2.EndeType, endeType)
	assert.Equal(t, planData2.PublicKey, publicABytes)
	assert.Equal(t, false, planData2.EnData)
}
