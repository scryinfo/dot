package denconfig

import (
	"crypto/ed25519"
	"crypto/sha512"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/crypto/curve25519"

	"github.com/scryinfo/scryg/sutils/sfile"
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
	flag.StringVar(&cfg.PrivateFileName, "private", "den.pri", "Ed25519 private file")
	flag.StringVar(&cfg.PubFileName, "public", "den.pub", "Ed25519 public file")
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

func GenEd25519(config *Config) error {
	{
		var err error
		config.PrivateFileName, err = filepath.Abs(config.PrivateFileName)
		if err != nil {
			slog.Error("", err)
			return err
		}
		config.PubFileName, err = filepath.Abs(config.PubFileName)
		if err != nil {
			slog.Error("", err)
			return err
		}
	}
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		slog.Error("generate key : ", err)
		return err
	}
	{
		dd := filepath.Dir(config.PrivateFileName)
		if dd != "" {
			if !sfile.ExistDir(dd) {
				err = os.MkdirAll(dd, 0o755)
				if err != nil {
					slog.Error("", err)
					return err
				}
			}
		}
	}
	err = os.WriteFile(config.PrivateFileName, priv, 0o600)
	if err != nil {
		slog.Error("", err)
		return err
	}
	slog.Info(fmt.Sprintf("private key: %s\n len: %d\n %#v", config.PrivateFileName, len(priv), priv))
	{
		dd := filepath.Dir(config.PubFileName)
		if dd != "" {
			if !sfile.ExistDir(dd) {
				err = os.MkdirAll(dd, 0o755)
				if err != nil {
					slog.Error("", err)
					return err
				}
			}
		}
	}
	err = os.WriteFile(config.PubFileName, pub, 0o600)
	if err != nil {
		slog.Error("", err)
		return err
	}
	slog.Info(fmt.Sprintf("public key: %s\n len: %d\n%#v", config.PubFileName, len(pub), pub))
	return nil
}

func Exit(err error) {
	code := 0
	if err != nil {
		code = 1
	}
	os.Exit(code)
}

func EdPrivToX25519Priv(edPriv ed25519.PrivateKey) ([]byte, error) {
	// ed25519 私钥结构：[32字节种子][32字节公钥]
	seed := edPriv[:32]
	h := sha512.New()
	h.Write(seed)
	digest := h.Sum(nil)
	// X25519 私钥规范裁剪
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64
	return digest[:32], nil
}

// Ed25519 公钥 → X25519 公钥
func EdPubToX25519Pub(edPub ed25519.PublicKey) ([]byte, error) {
	return curve25519.EdwardsToMontgomery(edPub)
}

func EncriptionConfig(config []byte, pri ed25519.PrivateKey) (error, []byte) {
	newPub, newPri, err := ed25519.GenerateKey(nil)
	if err != nil {
		slog.Error("generate key : ", err)
		return err, nil
	}

}
