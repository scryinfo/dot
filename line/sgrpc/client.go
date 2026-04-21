package sgrpc

type RpcClient struct {
	config RpcClientConfig
}

func NewRpcClient(config RpcClientConfig) (*RpcClient, func(), error) {
	d := &RpcClient{config: config}

	return d, d.stop, nil
}

func (p *RpcClient) stop() {

}
