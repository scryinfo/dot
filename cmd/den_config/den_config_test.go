package denconfig

import (
	"crypto/ecdh"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeencrypt(t *testing.T) {
	a := assert.New(t)
	priv, err := ecdh.X25519().GenerateKey(nil)
	a.Nil(err)
	data := []byte("testsdk sdfdsf")
	prepub, err := Encription(data, priv.PublicKey())
	a.Nil(err)
	decrypted, err := Decription(prepub, priv)
	a.Nil(err)
	a.Equal(data, decrypted)

}
