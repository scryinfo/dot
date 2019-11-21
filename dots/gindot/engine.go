// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gindot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

const (
	//EngineTypeId for gin dot
	EngineTypeId = "4943e959-7ad7-42c6-84dd-8b24e9ed30bb"
	//EngineLiveId for gin dot
	EngineLiveId = "4943e959-7ad7-42c6-84dd-8b24e9ed30bb"
)

type configEngine struct {
	Addr         string   `json:"addr"`         // addr smaple:  ":8080"
	KeyFile      string   `json:"keyFile"`      //if it is not abs path, preferred to use the executable path
	PemFile      string   `json:"pemFile"`      //if it is not abs path, preferred to use the executable path
	LogSkipPaths []string `json:"logSkipPaths"` // not write info log, sample: ["/tt", "/other"]
}

//GinEngine  gin dot
type Engine struct {
	ginEngine     *gin.Engine
	config        configEngine
	loggerOnlyGin dot.SLogger
}

//DefaultGinEngine return the default gin dot,
//it have to call after the line ceated
func DefaultGinEngine() *gin.Engine {
	logger := dot.Logger()
	l := dot.GetDefaultLine()
	if l == nil {
		logger.Errorln("the line do not create, do not call it")
		return nil
	}
	d, err := l.ToInjecter().GetByLiveId(EngineLiveId)
	if err != nil {
		logger.Errorln(err.Error())
		return nil
	}

	if g, ok := d.(*Engine); ok {
		return g.ginEngine
	}

	logger.Errorln("do not get the gin dot")
	return nil
}

//construct dot
func newGinDot(conf []byte) (dot.Dot, error) {
	dconf := &configEngine{}
	err := dot.UnMarshalConfig(conf, dconf)
	if err != nil {
		return nil, err
	}

	d := &Engine{config: *dconf}

	return d, err
}

//TypeLiveGinDot generate data for structural  dot
func TypeLiveGinDot() *dot.TypeLives {
	return &dot.TypeLives{
		Meta: dot.Metadata{TypeId: EngineTypeId, NewDoter: func(conf []byte) (dot dot.Dot, err error) {
			return newGinDot(conf)
		}},
	}
}

//jayce edit
//return config of GinDot
func ConfigTypeLiveGinDot() *dot.ConfigTypeLives {
	paths := make([]string, 0)
	paths = append(paths, "")
	return &dot.ConfigTypeLives{
		TypeIdConfig: EngineTypeId,
		ConfigInfo: &configEngine{
			LogSkipPaths: paths,
		},
	}
}

//Create create the gin
func (c *Engine) Create(l dot.Line) error {
	c.ginEngine = gin.New()
	c.loggerOnlyGin = dot.Logger().NewLogger(1)
	c.ginEngine.Use(c.makeLogger(l), gin.Recovery())
	return nil
}

//AfterAllStart run the function after start
func (c *Engine) AfterAllStart(l dot.Line) {
	go c.startServer()
}

func (c *Engine) GinEngine() *gin.Engine {
	return c.ginEngine
}

//all post
func (c *Engine) RouterPost(h interface{}, pre string) {
	post := reflect.ValueOf(c.ginEngine).MethodByName("POST")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		post.Call(vs)
	})
}

//all get
func (c *Engine) RouterGet(h interface{}, pre string) {
	get := reflect.ValueOf(c.ginEngine).MethodByName("GET")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		get.Call(vs)
	})
}

func (c *Engine) startServer() {
	llog := dot.Logger() //do not use the c.loggerOnlyGin, it only for gin
	if len(c.config.KeyFile) > 0 && len(c.config.PemFile) > 0 {
		keyFile := c.config.KeyFile
		pemFile := c.config.PemFile
		{
			ex, err := os.Executable()
			if err == nil {
				ex = filepath.Dir(ex)
			} else {
				ex = ""
			}
			if !filepath.IsAbs(keyFile) { //preferred to use the executable path
				t := filepath.Join(ex, keyFile)
				if sfile.ExistFile(t) {
					keyFile = t
				}
			}

			if !filepath.IsAbs(pemFile) { //preferred to use the executable path
				t := filepath.Join(ex, pemFile)
				if sfile.ExistFile(t) {
					pemFile = t
				}
			}
			if sfile.ExistFile(pemFile) && sfile.ExistFile(keyFile) {
				err := c.ginEngine.RunTLS(c.config.Addr, pemFile, keyFile)
				if err != nil {
					llog.Errorln(err.Error())
				}
			} else {
				llog.Errorln("the keyfile or pemfile do not exist")
				return
			}
		}
	} else {
		err := c.ginEngine.Run(c.config.Addr)
		if err != nil {
			llog.Errorln(err.Error())
		}
	}
}

func (c *Engine) makeLogger(l dot.Line) gin.HandlerFunc {

	formatter := defaultLogFormatter
	notLogged := c.config.LogSkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}
	logger := c.loggerOnlyGin

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			logger.Info(func() string {
				return formatter(param)
			})
		}
	}
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %v |%3d %s| %13v | %15s |%s %-7s %s %s\n%s",
		//param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
