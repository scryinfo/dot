package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

type handle struct {
	service interface{}
	preName string
	handle  *jsonrpc.Server
}

func (h *handle) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	h.handle.ServeHTTP(res, req)
}

func NewHandle(preName string, service interface{}) http.Handler {
	h := &handle{
		service: service,
		preName: preName,
	}
	h.handle = jsonrpc.NewServer(makeJsonrpc(preName, service), jsonrpc.ServerErrorLogger(&rpcLogger{}))
	return h
}

func NewHandles(servives map[string]interface{}) http.Handler {
	methods := jsonrpc.EndpointCodecMap{}
	for k, v := range servives {
		temp := makeJsonrpc(k, v)
		for tk, tv := range temp {
			methods[tk] = tv
		}
	}
	h := jsonrpc.NewServer(methods, jsonrpc.ServerErrorLogger(&rpcLogger{}))
	return h
}

// json rpc的解码函数
func nopDecoder(ctx context.Context, j json.RawMessage) (interface{}, error) {
	return j, nil
}

// json rpc的编码函数
func nopEncoder(ctx context.Context, req interface{}) (json.RawMessage, error) {
	bs, err := json.Marshal(req)
	return bs, err
}

func makeJsonrpc(preName string, server interface{}) jsonrpc.EndpointCodecMap {
	ecm := jsonrpc.EndpointCodecMap{}
	receiver := reflect.ValueOf(server)
	typ := receiver.Type()
	for i := 0; i < receiver.NumMethod(); i++ {
		tmethod := typ.Method(i)
		vmethod := receiver.Method(i)
		if tmethod.PkgPath != "" {
			continue // tmethod not exported
		}
		fn := tmethod.Type
		if fn.NumOut() == 1 && fn.NumIn() == 2 { //只处理单个参数
			endpointCodec := jsonrpc.EndpointCodec{
				Decode: nopDecoder,
				Encode: nopEncoder,
			}
			endpointCodec.Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
				defer func() {
					if e2 := recover(); e2 != nil {
						response = nil
						err = nil //todo 不太确定 err的具体工作过程， 这里暂定
						dot.Logger().Debugln("ApiService", zap.Any("", e2))
					}
				}()
				inType := fn.In(1)
				inValue := reflect.New(inType)
				bs, ok := request.(json.RawMessage)
				if ok && !bytes.Equal(bs, []byte("[null]")) { //"[null]", 就是 nil
					err = json.Unmarshal(bs, inValue.Interface())
					if err != nil {
						dot.Logger().Errorln("ApiService", zap.Error(err))
						inValue = reflect.New(inType)
					}
				} else {
					inValue = reflect.New(inType)
				}
				rev := vmethod.Call([]reflect.Value{inValue.Elem()})

				response = rev[0].Interface()
				err = nil
				return
			}
			if len(preName) > 0 {
				ecm[preName+tmethod.Name] = endpointCodec
			} else {
				ecm[tmethod.Name] = endpointCodec
			}
		}
	}
	return ecm
}

type rpcLogger struct{}

func (rpcLogger) Log(params ...interface{}) error {
	var err interface{} = params
	if len(params) == 2 {
		err = params[1]
	}
	dot.Logger().Warnln("ApiService rpc", zap.Any("", err))
	return nil
}
