package scrypto

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestX25519_ComputeSecret(t *testing.T) {

	ecdh := X25519()

	secret := make([]byte, 32)
	for i := 0; i < 10; i++ {
		privateA, publicA, err := ecdh.GenerateKey(nil)
		assert.Equal(t, nil, err, "A key pair generation failed")

		privateB, publicB, err := ecdh.GenerateKey(nil)
		assert.Equal(t, nil, err, "B key pair generation failed")

		secA, err := ecdh.ComputeSecret(privateA, publicB)
		assert.Equal(t, nil, err, "A compute error")
		secB, err := ecdh.ComputeSecret(privateB, publicA)
		assert.Equal(t, nil, err, "B compute error")

		assert.Equal(t, secA, secB, fmt.Sprintf("DH failed: secrets are not equal:\nA got: %s\nB got: %s", hex.EncodeToString(secA), hex.EncodeToString(secB)))
		assert.NotEqual(t, secret, secA, "DH generates the same secret all the time")
		copy(secret, secA)
	}

}
