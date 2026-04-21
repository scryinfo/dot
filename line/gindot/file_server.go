package gindot

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/gddo/httputil/header"
)

const (
	gzip = "gzip"
	br   = "br"
	con
)

// see https://github.com/lpar/gzipped/blob/master/fileserver.go
// FileServer is a drop-in replacement for Go's standard http.FileServer
// which adds support for static resources precompressed with gzip, at
// the cost of removing the support for directory browsing.
//
// If file filename.ext has a compressed version filename.ext.gz alongside
// it, if the client indicates that it accepts gzip-compressed data, and
// if the .gz file can be opened, then the compressed version of the file
// will be sent to the client. Otherwise the request is passed on to
// http.ServeContent, and the raw (uncompressed) version is used.
//
// It is up to you to ensure that the compressed and uncompressed versions
// of files match and have sensible timestamps.
//
// Compressed or not, requests are fulfilled using http.ServeContent, and
// details like accept ranges and content-type sniffing are handled by that
// method.
type FileServer struct {
	fs        http.FileSystem
	resPath   string
	paramName string
	encodings encodingByPreference
}

func NewFileServer(resPath string, paramName string, encodings []Encoding) *FileServer {
	fs := http.Dir(resPath)
	if len(encodings) < 1 {
		encodings = supportedEncodings[:]
	}
	return &FileServer{
		fs:        fs,
		resPath:   resPath,
		paramName: paramName,
		encodings: encodings,
	}
}

func (c *FileServer) Handler(ctx *gin.Context) {
	fileString := ctx.Param(c.paramName)
	fileString = path.Clean(fileString)
	if strings.HasSuffix(fileString, "/") { //不支持目录
		ctx.Status(http.StatusNotFound)
		ctx.Abort()
		return
	}

	// Find the best acceptable file, including trying uncompressed
	if file, info, err := c.findBestFile(ctx.Writer, ctx.Request, fileString); err == nil {
		http.ServeContent(ctx.Writer, ctx.Request, fileString, info.ModTime(), file)
		file.Close()
		//ctx.Abort()
		return
	}

	ctx.Status(http.StatusNotFound)
	ctx.Abort()
	// return
}

// Encoding represents an Accept-Encoding. All of these fields are pre-populated
// in the supportedEncodings variable, except the ClientPreference which is updated
// (by copying a value from supportedEncodings) when examining client headers.
type Encoding struct {
	Name             string  `json:"name"`             // the Encoding Name
	Extension        string  `json:"extension"`        // the file Extension (including a leading dot)
	ClientPreference float64 `json:"clientPreference"` // the client's preference
	ServerPreference int     `json:"serverPreference"` // the server's preference
}

// Helper type to sort encodings, using ClientPreference first, and then
// ServerPreference as a tie breaker. This sorts in *DESCENDING* order, rather
// than the usual ascending order.
type encodingByPreference []Encoding

// Implement the sort.Interface interface
func (e encodingByPreference) Len() int { return len(e) }
func (e encodingByPreference) Less(i, j int) bool {
	if e[i].ClientPreference == e[j].ClientPreference {
		return e[i].ServerPreference > e[j].ServerPreference
	}
	return e[i].ClientPreference > e[j].ClientPreference
}
func (e encodingByPreference) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// Supported encodings. Higher server preference means the Encoding will be when
// the client doesn't have an explicit preference.
var supportedEncodings = [...]Encoding{
	{
		Name:             "gzip",
		Extension:        ".gz",
		ServerPreference: 1,
	},
	{
		Name:             "br",
		Extension:        ".br",
		ServerPreference: 2,
	},
}

func (c *FileServer) openAndStat(path string) (http.File, os.FileInfo, error) {
	file, err := c.fs.Open(path)
	var info os.FileInfo
	// This slightly weird variable reuse is so we can get 100% test coverage
	// without having to come up with a test file that can be opened, yet
	// fails to stat.
	if err == nil {
		info, err = file.Stat()
	}
	if err != nil {
		return file, nil, err
	}
	if info.IsDir() {
		return file, nil, fmt.Errorf("%s is directory", path)
	}
	return file, info, nil
}

// Build a []Encoding based on the Accept-Encoding header supplied by the
// client. The returned list will be sorted from most-preferred to
// least-preferred.
func (c *FileServer) acceptable(r *http.Request) []Encoding {
	// list of acceptable encodings, as provided by the client
	acceptEncodings := make([]Encoding, 0, len(c.encodings))

	// the quality of the * Encoding; this will be -1 if not sent by client
	starQuality := -1.

	// encodings we've already seen (used to handle duplicates and *)
	seenEncodings := make(map[string]interface{})

	// match the client accept encodings against the ones we support
	for _, aspec := range header.ParseAccept(r.Header, acceptEncodingHeader) {
		if _, alreadySeen := seenEncodings[aspec.Value]; alreadySeen {
			continue
		}
		seenEncodings[aspec.Value] = nil
		if aspec.Value == "*" {
			starQuality = aspec.Q
			continue
		}
		for _, known := range c.encodings {
			if aspec.Value == known.Name && aspec.Q != 0 {
				enc := known
				enc.ClientPreference = aspec.Q
				acceptEncodings = append(acceptEncodings, enc)
				break
			}
		}
	}

	// If the client sent Accept: *, add all our extra known encodings. Use
	// the quality of * as the client quality for the Encoding.
	if starQuality != -1. {
		for _, known := range c.encodings {
			if _, seen := seenEncodings[known.Name]; !seen {
				enc := known
				enc.ClientPreference = starQuality
				acceptEncodings = append(acceptEncodings, enc)
			}
		}
	}

	// sort the Encoding based on client/server preference
	sort.Sort(encodingByPreference(acceptEncodings))
	return acceptEncodings
}

const (
	acceptEncodingHeader  = "Accept-Encoding"
	contentEncodingHeader = "Content-Encoding"
	contentLengthHeader   = "Content-Length"
	rangeHeader           = "Range"
	varyHeader            = "Vary"
)

// Find the best file to serve based on the client's Accept-Encoding, and which
// files actually exist on the filesystem. If no file was found that can satisfy
// the request, the error field will be non-nil.
func (c *FileServer) findBestFile(w http.ResponseWriter, r *http.Request, fpath string) (http.File, os.FileInfo, error) {
	// find the best matching file
	for _, enc := range c.acceptable(r) {
		if file, info, err := c.openAndStat(fpath + enc.Extension); err == nil {
			wHeader := w.Header()
			wHeader[contentEncodingHeader] = []string{enc.Name}
			wHeader.Add(varyHeader, acceptEncodingHeader)

			if len(r.Header[rangeHeader]) == 0 {
				// If not a range request then we can easily set the content length which the
				// Go standard library does not do if "Content-Encoding" is set.
				wHeader[contentLengthHeader] = []string{strconv.FormatInt(info.Size(), 10)}
			}
			return file, info, nil
		}
	}

	// if nothing found, try the base file with no content-Encoding
	return c.openAndStat(fpath)
}

//func (c *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	upath := r.URL.Path
//	if !strings.HasPrefix(upath, "/") {
//		upath = "/" + upath
//		r.URL.Path = upath
//	}
//	fpath := path.Clean(upath)
//	if strings.HasSuffix(fpath, "/") {
//		// If you wanted to put back directory browsing support, this is
//		// where you'd do it.
//		http.NotFound(w, r)
//		return
//	}
//
//	// Find the best acceptable file, including trying uncompressed
//	if file, info, err := c.findBestFile(w, r, fpath); err == nil {
//		http.ServeContent(w, r, fpath, info.ModTime(), file)
//		file.Close()
//		return
//	}
//
//	// Doesn't exist, compressed or uncompressed
//	http.NotFound(w, r)
//}
