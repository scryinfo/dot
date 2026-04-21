package sgrpc

type RpcServer struct {
	config RpcServerConfig
}

func NewRpcServer(config RpcServerConfig) (*RpcServer, func(), error) {
	d := &RpcServer{config: config}

	return d, d.stop, nil
}

func (s *RpcServer) stop() {

}
