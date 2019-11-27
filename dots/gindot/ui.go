package gindot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

const (
	UiTypeId = "d9972be7-cef9-464c-9bbb-d1f11abea803"
)

type configUi struct {
	UrlRelativePath string            `json:"urlRelativePath"` //url 的相对路径
	ResRelativePath string            `json:"resRelativePath"` //资源的相对路径， 查找的先后为： 绝对路路, 相对路径，执行路径，当前路径，当前用户所在路径，当都没有找到时认为文件没有找到
	Paths           map[string]string `json:"paths"`           //静态资源的路径，可以是文件或目录, key: pre, value: folder or file path
}

//Ui  用于静态资源的组件
type Ui struct {
	Engine_ *Engine `dot:""`

	router       *gin.RouterGroup
	config       configUi
	executePath  string //执行文件所在路径
	currentPath  string //当前路或工作路径
	userPath     string //当前用户所在路径
	relativePath string //相对路径，当来自于配置文件
}

func (c *Ui) AfterAllInject(l dot.Line) {
	c.router = c.Engine_.GinEngine().Group(c.config.UrlRelativePath)
}

//Start start the gin
func (c *Ui) Start(ignore bool) error {
	logger := dot.Logger()
	for k, v := range c.config.Paths {
		res := c.ResAbsolutePath(v)
		if len(res) > 0 {
			if sfile.IsDir(res) {
				c.router.Static(k, res)
			} else if sfile.IsFile(res) {
				c.router.StaticFile(k, res)
			} else {
				logger.Errorln("", zap.String("", "can not: "+v+" realy: "+res))
			}
		} else {
			logger.Errorln("", zap.String("", fmt.Sprintf("can not find : %s  under (%s, %s, %s, %s)", v, c.relativePath, c.executePath, c.currentPath, c.userPath)))
		}
	}
	return nil
}

func (c *Ui) Router() *gin.RouterGroup {
	return c.router
}

func (c *Ui) UrlRelativePath() string {
	return c.config.UrlRelativePath
}

func (c *Ui) ResRelativePath() string {
	return c.relativePath
}

func (c *Ui) SetResRelativePath(relativePath string) {
	c.relativePath = ""
	c.relativePath = c.ResAbsolutePath(relativePath)
}

//查找的先后为： 绝对路路, 相对路径，执行路径，当前路径，当前用户所在路径，当都没有找到时认为文件没有找到
//没有找到文件时，返回""
func (c *Ui) ResAbsolutePath(res string) string {
	if filepath.IsAbs(res) {
		if sfile.ExistFile(res) {
			return res
		}
	}

	if len(c.relativePath) > 0 {
		temp := filepath.Join(c.relativePath, res)
		if sfile.ExistFile(temp) {
			return temp
		}
	}

	if len(c.executePath) > 0 {
		temp := filepath.Join(c.executePath, res)
		if sfile.ExistFile(temp) {
			return temp
		}
	}

	if len(c.currentPath) > 0 {
		temp := filepath.Join(c.currentPath, res)
		if sfile.ExistFile(temp) {
			return temp
		}
	}

	if len(c.userPath) > 0 {
		temp := filepath.Join(c.userPath, res)
		if sfile.ExistFile(temp) {
			return temp
		}
	}

	return ""
}

//construct dot
func newUi(conf []byte) (*Ui, error) {
	dconf := &configUi{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	ui := &Ui{config: *dconf}
	d := ui
	{
		ui.executePath, err = os.Executable()
		if err != nil {
			dot.Logger().Errorln("Ui", zap.Error(err))
			ui.executePath = ""
		} else {
			ui.executePath = filepath.Dir(ui.executePath)
		}
		ui.currentPath, err = os.Getwd()
		if err != nil {
			dot.Logger().Errorln("Ui", zap.Error(err))
			ui.currentPath = ""
		}
		ui.userPath, err = os.UserHomeDir()
		if err != nil {
			dot.Logger().Errorln("Ui", zap.Error(err))
			ui.userPath = ""
		}

		if len(ui.config.ResRelativePath) > 0 {
			ui.relativePath = ui.ResAbsolutePath(ui.config.ResRelativePath)
		}

		err = nil
	}

	return d, err
}

//UiTypeLives generate data for structural  dot,  include gindot.Engine
func UiTypeLives() []*dot.TypeLives {
	return []*dot.TypeLives{&dot.TypeLives{
		Meta: dot.Metadata{TypeId: UiTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newUi(conf)
		}},
	},
		TypeLiveGinDot(),
	}
}

//return config of Ui
func UiConfigTypeLive() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: UiTypeId,
		ConfigInfo:   &configUi{},
	}
}
