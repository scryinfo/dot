package denconfig

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/scryinfo/scryg/sutils/sfile"
	"golang.org/x/crypto/hkdf"
)

type Config struct {
	PrivateFileName  string
	PubFileName      string
	ConfigFileName   string
	DeConfigFileName string
	EnConfigFileName string
}

func ParseFlags() Config {
	var cfg Config
	flag.StringVar(&cfg.PrivateFileName, "private", "den_x25519.pri", "x25519 private file")
	flag.StringVar(&cfg.PubFileName, "public", "den_x25519.pub", "x25519 public file")
	flag.StringVar(&cfg.ConfigFileName, "config", "config.toml", "config file")
	flag.Parse()
	return cfg
}

func (p *Config) Abs() error {
	var err error
	p.ConfigFileName, err = filepath.Abs(p.ConfigFileName)
	if err != nil {
		return err
	}
	p.PrivateFileName, err = filepath.Abs(p.PrivateFileName)
	if err != nil {
		return err
	}
	p.PubFileName, err = filepath.Abs(p.PubFileName)
	if err != nil {
		return err
	}
	return nil
}

func (p *Config) ConfigFile() error {
	if p.ConfigFileName == "" {
		return fmt.Errorf("the config file dont exist: %s", p.ConfigFileName)
	}
	dir := filepath.Dir(p.ConfigFileName)
	ext := filepath.Ext(p.ConfigFileName)
	name := filepath.Base(p.ConfigFileName)
	name = name[:len(name)-len(ext)]
	p.DeConfigFileName = filepath.Join(dir, fmt.Sprintf("%s_de%s", name, ext))
	p.EnConfigFileName = filepath.Join(dir, fmt.Sprintf("%s_en%s", name, ext))
	return nil
}

func GenX25519(config *Config) error {
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
func Encription(config []byte, pub *ecdh.PublicKey) ([]byte, error) {
	priv, err := ecdh.X25519().GenerateKey(nil)
	if err != nil {
		slog.Error("generate key : ", "", err)
		return nil, err
	}
	sharedSecret, err := priv.ECDH(pub)
	if err != nil {
		slog.Error("ecdh : ", "", err)
		return nil, err
	}
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

	cipherText := aesGCM.Seal(nonce, nonce, config, nil)
	prepub := make([]byte, 0, len(cipherText)+len(pub.Bytes()))
	prepub = append(prepub, priv.PublicKey().Bytes()...)
	prepub = append(prepub, cipherText...)
	return prepub, nil
}
func Decription(prepub []byte, priv *ecdh.PrivateKey) ([]byte, error) {
	if len(prepub) < 32 {
		return nil, fmt.Errorf("the prepub is too short")
	}

	pub, err := ecdh.X25519().NewPublicKey(prepub[:32])
	if err != nil {
		return nil, err
	}
	sharedSecret, err := priv.ECDH(pub)
	if err != nil {
		return nil, err
	}
	cipherText := prepub[32:]
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
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("the cipher text is too short")
	}
	nonce, actualCipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
