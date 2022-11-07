package gindot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const UiTypeID = "d9972be7-cef9-464c-9bbb-d1f11abea803" //type id of dot

type configUi struct {
	UrlRelativePath string `json:"urlRelativePath"` //relative path of url
	ResRelativePath string `json:"resRelativePath"` //relative path of resource， The order of locating files is: absolute path, relative path，executable path，current path，user path
	Paths           []struct {
		RelativePath string `json:"relativePath"`
		Value        string `json:"value"` // the static resource of path, it is file or folder
	} `json:"paths"`
	MainHTMlName string     `json:"mainHTMlName"` //maybe home.html or index.html
	NoCompress   bool       `json:"noCompress"`
	Encodings    []Encoding `json:"encodings"` //default are br-->gzip
}

//Ui  add static resource into gin
type Ui struct {
	Engine_ *Engine `dot:""`

	router       *gin.RouterGroup
	config       configUi
	executePath  string //executable path
	currentPath  string //current path
	userPath     string //user path
	relativePath string //relative path
}

func (c *Ui) Injected(l dot.Line) error {
	c.router = c.Engine_.GinEngine().Group(c.config.UrlRelativePath)
	c.router.Use(func(ctx *gin.Context) {
		ctx.Handler()
	})
	if len(c.config.MainHTMlName) > 0 {
		c.Engine_.GinEngine().GET("/", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusFound, c.UrlRelativePath()+"/"+c.config.MainHTMlName)
		})
		c.Engine_.GinEngine().GET("/home", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusFound, c.UrlRelativePath()+"/"+c.config.MainHTMlName)
		})
		c.Engine_.GinEngine().GET("/index", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusFound, c.UrlRelativePath()+"/"+c.config.MainHTMlName)
		})
	}
	return nil
}

//Start start the gin
func (c *Ui) Start(ignore bool) error {
	logger := dot.Logger()
	for _, it := range c.config.Paths {
		res := c.ResAbsolutePath(it.Value)
		if len(res) > 0 {
			logger.Debugln("Ui", zap.String("", res))
			if sfile.IsDir(res) {
				if c.config.NoCompress {
					c.router.Static(it.RelativePath, res)
				} else {
					//urlPrePath := path.Join(c.router.BasePath(), it.RelativePath)
					//if it.RelativePath[len(it.RelativePath) -1] == '/' && urlPrePath[len(urlPrePath) -1] != '/' {
					//	urlPrePath += "/"
					//}
					if len(c.config.Encodings) < 1 {
						c.config.Encodings = supportedEncodings[:]
					}
					handler := NewFileServer(res, "filepath", c.config.Encodings)
					urlPattern := path.Join(it.RelativePath, "/*filepath")
					c.router.GET(urlPattern, handler.Handler)
					c.router.HEAD(urlPattern, handler.Handler)
				}
			} else if sfile.IsFile(res) {
				c.router.StaticFile(it.RelativePath, res)
			} else {
				logger.Errorln("", zap.String("", "can not: "+it.Value+" realy: "+res))
			}
		} else {
			logger.Errorln("", zap.String("", fmt.Sprintf("can not find : %s  under (%s, %s, %s, %s)", it.RelativePath, c.relativePath, c.executePath, c.currentPath, c.userPath)))
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

//ResAbsolutePath the order of locating files is: absolute path, relative path，executable path，current path，user path
//if do not find, then return ""
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
			////for dev
			if !sfile.ExistFile(ui.ResAbsolutePath("dist")) {
				dot.Logger().Debugln("UiPreAdd", zap.String("", ui.config.ResRelativePath))
				ui.SetResRelativePath(ui.config.ResRelativePath)
			}
		}

		err = nil
	}

	return d, err
}

//UiTypeLives generate data for structural  dot,  include gindot.Engine
func UiTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: UiTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newUi(conf)
		}},

		Lives: []dot.Live{
			{
				LiveID:    UiTypeID,
				RelyLives: map[string]dot.LiveID{"Engine_": EngineLiveID},
			},
		},
	}}
	lives = append(lives, GinDotTypeLives()...)
	return lives
}

//return config of Ui
func UiConfigTypeLive() *dot.ConfigTypeLive {
	return &dot.ConfigTypeLive{
		TypeIDConfig: UiTypeID,
		ConfigInfo:   &configUi{},
	}
}
