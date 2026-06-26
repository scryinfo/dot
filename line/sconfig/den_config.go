package sconfig

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/hkdf"
)

const (
	salt = "_dot_solt&*_"
	info = "_dot_info*&_"
)

func deriveKey(sharedSecret []byte) ([]byte, error) {
	hkdfReader := hkdf.New(sha256.New, sharedSecret, []byte(salt), []byte(info))
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return nil, err
	}
	return aesKey, nil
}
func EncriptionFile(plainConfig []byte, pub *ecdh.PublicKey) ([]byte, error) {
	priv, err := ecdh.X25519().GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := priv.ECDH(pub)
	if err != nil {
		return nil, err
	}
	cipherText, err := EncriptionKey(plainConfig, sharedSecret)
	if err != nil {
		return nil, err
	}

	prepub := make([]byte, 0, len(cipherText)+len(pub.Bytes()))
	prepub = append(prepub, priv.PublicKey().Bytes()...)
	prepub = append(prepub, cipherText...)
	return prepub, nil
}
func DecriptionFile(prepubConfig []byte, priv *ecdh.PrivateKey) ([]byte, error) {
	if len(prepubConfig) < 32 {
		return nil, fmt.Errorf("the prepub is too short")
	}

	pub, err := ecdh.X25519().NewPublicKey(prepubConfig[:32])
	if err != nil {
		return nil, err
	}
	sharedSecret, err := priv.ECDH(pub)
	if err != nil {
		return nil, err
	}
	cipherText := prepubConfig[32:]

	plainText, err := DecriptionKey(cipherText, sharedSecret)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func EncriptionKey(plainValue []byte, sharedSecret []byte) ([]byte, error) {
	key, err := deriveKey(sharedSecret)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainValue, nil)
	return cipherText, nil
}
func DecriptionKey(cipherValue []byte, sharedSecret []byte) ([]byte, error) {
	if len(cipherValue) < 32 {
		return nil, fmt.Errorf("the cipherText is too short")
	}

	key, err := deriveKey(sharedSecret)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(cipherValue) < nonceSize {
		return nil, fmt.Errorf("the cipher text is too short")
	}
	nonce, actualCipherText := cipherValue[:nonceSize], cipherValue[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func ListCdConfigFiles() ([]string, error) {
	var confs []string
	var exeDir string
	{
		exeFile, err := os.Executable()
		if err != nil {
			return nil, err
		}
		exeDir = filepath.Dir(exeFile)
	}
	entries, err := os.ReadDir(exeDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			ex := filepath.Ext(entry.Name())
			if ex == ExtensionNameToml || ex == ExtensionNameJson || ex == ExtensionNameYaml {
				confs = append(confs, entry.Name())
			}
		}
	}
	return confs, nil
}
