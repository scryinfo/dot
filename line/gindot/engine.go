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
	"github.com/rs/zerolog"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

type EngineConfig struct {
	Addr         string   `json:"addr" yaml:"addr"`                 // addr smaple:  ":8080"
	KeyFile      string   `json:"keyFile" yaml:"keyFile"`           //if it is not abs path, preferred to use the executable path
	PemFile      string   `json:"pemFile" yaml:"pemFile"`           //if it is not abs path, preferred to use the executable path
	LogSkipPaths []string `json:"logSkipPaths" yaml:"logSkipPaths"` // not write info log, sample: ["/tt", "/other"]
}

// GinEngine  gin dot
type Engine struct {
	ginEngine     *gin.Engine
	config        EngineConfig
	loggerOnlyGin *dot.LoggerType
}

// construct dot
func NewGinDot(conf *EngineConfig, loggerOnlyGin *dot.LoggerType) (*Engine, error) {
	d := &Engine{config: *conf, ginEngine: gin.New(), loggerOnlyGin: loggerOnlyGin}
	d.ginEngine.Use(d.makeLogger(), gin.Recovery())
	go d.startServer()
	return d, nil
}

func (c *Engine) GinEngine() *gin.Engine {
	return c.ginEngine
}

// all post
func (c *Engine) RouterPost(h interface{}, pre string) {
	post := reflect.ValueOf(c.ginEngine).MethodByName("POST")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		post.Call(vs)
	})
}

// all get
func (c *Engine) RouterGet(h interface{}, pre string) {
	get := reflect.ValueOf(c.ginEngine).MethodByName("GET")
	RouterSelf(h, pre, func(url string, gmethod reflect.Value) {
		vs := []reflect.Value{reflect.ValueOf(url), gmethod}
		get.Call(vs)
	})
}

func (c *Engine) startServer() {
	logger := &dot.Logger //do not use the c.loggerOnlyGin, it only for gin
	if logger.GetLevel() != zerolog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}
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
					logger.Error().Err(err).Send()
				}
			} else {
				logger.Error().Msg("the keyfile or pemfile do not exist")
				return
			}
		}
	} else {
		err := c.ginEngine.Run(c.config.Addr)
		if err != nil {
			logger.Error().Err(err).Send()
		}
	}
}

func (c *Engine) makeLogger() gin.HandlerFunc {

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
			logger.Info().Msg(formatter(param))
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
