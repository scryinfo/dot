// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scryinfo/dot/dot"
	"github.com/stretchr/testify/assert"
)

func TestEd25519_GenerateECDSAKey(t *testing.T) {
	ec := NewEd25519(dot.NewLogger(&dot.LogConfig{}))

	rootKey, err := MakeEd25519Key()
	assert.Nil(t, err)

	keyFile := "root.key"
	certFile := "root.cert"
	{
		exPath, err := os.Executable()
		assert.Nil(t, err)
		exPath = filepath.Dir(exPath)
		keyFile = filepath.Join(exPath, keyFile)
		certFile = filepath.Join(exPath, certFile)
	}
	_, err = ec.GenerateRoot(rootKey, keyFile, certFile, []string{"scry"}, []string{"scry"})

	defer func() {
		_ = os.Remove(keyFile)
		_ = os.Remove(certFile)
	}()
	assert.Nil(t, err)

	loadKey, err := ec.PrivateKey(keyFile)
	assert.Nil(t, err)
	assert.True(t, rootKey.Equal(loadKey))

	caPub, err := ec.PublicKey(certFile)
	assert.Nil(t, err)
	assert.True(t, caPub.Equal(rootKey.Public()))

}

func TestEd25519_GenerateCertKey(t *testing.T) {

	ec := NewEd25519(dot.NewLogger(&dot.LogConfig{}))

	rootKey, err := MakeEd25519Key()
	assert.Nil(t, err)

	rootKeyFile := "root.key"
	rootCertFile := "cert.cert"
	leafKeyFile := "leaf.key"
	leafCertFile := "leaf.cert"
	{
		exPath, _ := os.Executable()
		exPath = filepath.Dir(exPath)
		rootKeyFile = filepath.Join(exPath, rootKeyFile)
		rootCertFile = filepath.Join(exPath, rootCertFile)
		leafKeyFile = filepath.Join(exPath, leafKeyFile)
		leafCertFile = filepath.Join(exPath, leafCertFile)
	}
	rootCert, err := ec.GenerateRoot(rootKey, rootKeyFile, rootCertFile, []string{"scry"}, []string{"scry"})

	defer func() {
		_ = os.Remove(rootKeyFile)
		_ = os.Remove(rootCertFile)
		_ = os.Remove(leafKeyFile)
		_ = os.Remove(leafCertFile)
	}()

	assert.Nil(t, err)

	_, err = ec.GenerateLeaf(rootCert, rootKey, leafKeyFile, leafCertFile, []string{"scry"}, []string{"scry"})
	defer func() {
		_ = os.Remove(leafKeyFile)
		_ = os.Remove(leafCertFile)
	}()
	assert.Nil(t, err)

	leafKey, err := ec.PrivateKey(leafKeyFile)
	assert.Nil(t, err)

	caPub, err := ec.PublicKey(leafCertFile)
	assert.Nil(t, err)
	assert.True(t, caPub.Equal(leafKey.Public()))

}
