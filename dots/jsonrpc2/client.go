package jsonrpc2

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"net/url"
)

func CallOne(urlStr string, out interface{}, method string, params ...interface{}) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	client := jsonrpc.NewClient(u, method, jsonrpc.ClientResponseDecoder(func(_ context.Context, res jsonrpc.Response) (interface{}, error) {
		if res.Error != nil {
			return nil, *res.Error
		}
		err := json.Unmarshal(res.Result, out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}))
	_, err = client.Endpoint()(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}
