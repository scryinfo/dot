package denconfig

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
	"golang.org/x/crypto/hkdf"
)

func GenX25519(config *ConfigFile) error {
	{
		var err error
		config.PrivateFileName, err = filepath.Abs(config.PrivateFileName)
		if err != nil {
			slog.Error("", "", err)
			return err
		}
		config.PubFileName, err = filepath.Abs(config.PubFileName)
		if err != nil {
			slog.Error("", "", err)
			return err
		}
	}
	priv, err := ecdh.X25519().GenerateKey(nil)
	if err != nil {
		slog.Error("generate key : ", "", err)
		return err
	}
	{
		dd := filepath.Dir(config.PrivateFileName)
		if dd != "" {
			if !sfile.ExistDir(dd) {
				err = os.MkdirAll(dd, 0o755)
				if err != nil {
					slog.Error("", "", err)
					return err
				}
			}
		}
	}
	err = os.WriteFile(config.PrivateFileName, priv.Bytes(), 0o600)
	if err != nil {
		slog.Error("", "", err)
		return err
	}
	slog.Info(fmt.Sprintf("private key: %s\n len: %d\n %#v", config.PrivateFileName, len(priv.Bytes()), priv))
	{
		dd := filepath.Dir(config.PubFileName)
		if dd != "" {
			if !sfile.ExistDir(dd) {
				err = os.MkdirAll(dd, 0o755)
				if err != nil {
					slog.Error("", "", err)
					return err
				}
			}
		}
	}
	err = os.WriteFile(config.PubFileName, priv.PublicKey().Bytes(), 0o600)
	if err != nil {
		slog.Error("", "", err)
		return err
	}
	slog.Info(fmt.Sprintf("public key: %s\n len: %d\n%#v", config.PubFileName, len(priv.PublicKey().Bytes()), priv.PublicKey()))
	return nil
}

func Exit(err error) {
	code := 0
	if err != nil {
		code = 1
	}
	os.Exit(code)
}

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
		slog.Error("derive key : ", "", err)
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
			if ex == sconfig.ExtensionNameToml || ex == sconfig.ExtensionNameJson || ex == sconfig.ExtensionNameYaml {
				confs = append(confs, entry.Name())
			}
		}
	}
	return confs, nil
}
