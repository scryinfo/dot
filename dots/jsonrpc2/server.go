package jsonrpc2

import (
	"github.com/pkg/errors"
	"net/http"
)

//only one input parameter, one out parameter
func NewService(pattern string, preName string, service interface{}) (http.Handler, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}
	h := NewHandle(preName, service)
	return h, nil
}
