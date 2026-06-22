package denconfig

import (
	"crypto/ecdh"
	"testing"

	"github.com/scryinfo/dot/line/sconfig"
	"github.com/stretchr/testify/assert"
)

func TestDeencrypt(t *testing.T) {
	a := assert.New(t)
	priv, err := ecdh.X25519().GenerateKey(nil)
	a.Nil(err)
	data := []byte("testsdk sdfdsf")
	prepub, err := sconfig.EncriptionFile(data, priv.PublicKey())
	a.Nil(err)
	decrypted, err := sconfig.DecriptionFile(prepub, priv)
	a.Nil(err)
	a.Equal(data, decrypted)

}
