package imp

import (
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

	ciphertext, err := ende.EcdhEncode(privateA, publicB, plaintext)
	assert.Equal(t, nil, err)

	decodeText, err := ende.EcdhDecode(privateB, publicA, ciphertext)
	assert.Equal(t, nil, err)

	assert.Equal(t, plaintext, decodeText)
}
