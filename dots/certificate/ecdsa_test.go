// Scry Info.  All rights reserved.
// license that can be found in the license file.

package certificate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEcdsa_GenerateCaCertKey(t *testing.T) {
	ec := &Ecdsa{}

	caPri, err := MakePriKey()
	if err != nil {
		t.Error(err)
	}

	keyFile := "ca.key"
	pemFile := "ca.pem"
	{
		exPath, _ := os.Executable()
		exPath = filepath.Dir(exPath)
		keyFile = filepath.Join(exPath, keyFile)
		pemFile = filepath.Join(exPath, pemFile)
	}
	_, err = ec.GenerateCaCertKey(caPri, keyFile, pemFile, []string{"scry"}, []string{"scry"})

	defer func() {
		os.Remove(keyFile)
		os.Remove(pemFile)
	}()

	if err != nil {
		t.Error(err)
	}

	caPri2, err := ec.PrivateKey(keyFile)
	if err != nil {
		t.Error(err)
	}
	if err == nil {
		if caPri.D.Cmp(caPri2.D) != 0 {
			t.Error(err)
		}
	}

	caPub, err := ec.PublicKey(pemFile)
	if err != nil {
		t.Error(err)
	}
	if err == nil {
		if caPri.PublicKey.X.Cmp(caPub.X) != 0 || caPri.PublicKey.Y.Cmp(caPub.Y) != 0 {
			t.Error(err)
		}
	}

}

func TestEcdsa_GenerateCertKey(t *testing.T) {

	ec := &Ecdsa{}

	caPri, err := MakePriKey()
	if err != nil {
		t.Error(err)
	}

	keyFile := "ca.key"
	pemFile := "ca.pem"
	keySub := "sub.key"
	pemSub := "sum.pem"
	{
		exPath, _ := os.Executable()
		exPath = filepath.Dir(exPath)
		keyFile = filepath.Join(exPath, keyFile)
		pemFile = filepath.Join(exPath, pemFile)
		keySub = filepath.Join(exPath, keySub)
		pemSub = filepath.Join(exPath, pemSub)
	}
	ca, err := ec.GenerateCaCertKey(caPri, keyFile, pemFile, []string{"scry"}, []string{"scry"})

	defer func() {
		os.Remove(keyFile)
		os.Remove(pemFile)
	}()

	if err != nil {
		t.Error(err)
	}

	err = ec.GenerateCertKey(ca, caPri, keySub, pemSub, []string{"scry"}, []string{"scry"})
	defer func() {
		os.Remove(keySub)
		os.Remove(pemSub)
	}()
	if err != nil {
		t.Error(err)
	}

	subPri, err := ec.PrivateKey(keySub)
	if err != nil {
		t.Error(err)
	}

	caPub, err := ec.PublicKey(pemSub)
	if err != nil {
		t.Error(err)
	}
	if err == nil {
		if subPri.PublicKey.X.Cmp(caPub.X) != 0 || subPri.PublicKey.Y.Cmp(caPub.Y) != 0 {
			t.Error(err)
		}
	}

}
