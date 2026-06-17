package denconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scryinfo/scryg/sutils/sfile"
)

type ConfigKey struct {
	PrivateFileName  string
	PubFileName      string
	ConfigFileName   string
	DeConfigFileName string
	EnConfigFileName string
	ExeConfs         []string
	ExeDir           string
}

func (p *ConfigKey) Abs() error {
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

func (p *ConfigKey) ConfigFile() error {
	if len(p.ExeConfs) < 1 {
		err := p.InitExeConfs()
		if err != nil {
			return err
		}
	}
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

func (p *ConfigKey) InitExeConfs() error {
	cs, err := ListCdConfigFiles()
	if err != nil {
		return err
	}
	p.ExeConfs = cs
	exeFile, err := os.Executable()
	if err != nil {
		return err
	}
	p.ExeDir = filepath.Dir(exeFile)
	if p.ConfigFileName == "" || !sfile.ExistFile(filepath.Join(p.ExeDir, p.ConfigFileName)) {
		if len(p.ExeConfs) == 1 {
			p.ConfigFileName = p.ExeConfs[0]
		}
	}
	return nil
}
