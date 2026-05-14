// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scryinfo/dot/dot"
)

func TestEcdsa_GenerateECDSAKey(t *testing.T) {
	ec := NewEcdsa(dot.NewLogger(&dot.LogConfig{}))

	rootKey, err := MakeECDSAKey()
	if err != nil {
		t.Error(err)
	}

	keyFile := "root.key"
	certFile := "root.cert"
	{
		exPath, _ := os.Executable()
		exPath = filepath.Dir(exPath)
		keyFile = filepath.Join(exPath, keyFile)
		certFile = filepath.Join(exPath, certFile)
	}
	_, err = ec.GenerateRoot(rootKey, keyFile, certFile, []string{"scry"}, []string{"scry"})

	defer func() {
		_ = os.Remove(keyFile)
		_ = os.Remove(certFile)
	}()

	if err != nil {
		t.Error(err)
	}

	loadKey, err := ec.PrivateKey(keyFile)
	if err != nil {
		t.Error(err)
	}
	if err == nil {
		if !rootKey.Equal(loadKey) {
			t.Error(err)
		}
	}

	caPub, err := ec.PublicKey(certFile)
	if err != nil {
		t.Error(err)
	}
	if err == nil {
		if !rootKey.PublicKey.Equal(caPub) {
			t.Error(err)
		}
	}

}

func TestEcdsa_GenerateCertKey(t *testing.T) {

	ec := NewEcdsa(dot.NewLogger(&dot.LogConfig{}))

	rootKey, err := MakeECDSAKey()
	if err != nil {
		t.Error(err)
	}

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
	}()

	if err != nil {
		t.Error(err)
	}

	_, err = ec.GenerateLeaf(rootCert, rootKey, leafKeyFile, leafCertFile, []string{"scry"}, []string{"scry"})
	defer func() {
		_ = os.Remove(leafKeyFile)
		_ = os.Remove(leafCertFile)
	}()
	if err != nil {
		t.Error(err)
	}

	leafKey, err := ec.PrivateKey(leafKeyFile)
	if err != nil {
		t.Error(err)
	}

	caPub, err := ec.PublicKey(leafCertFile)
	if err != nil {
		t.Error(err)
	}
	if err == nil {
		if !leafKey.PublicKey.Equal(caPub) {
			t.Error(err)
		}
	}

}
