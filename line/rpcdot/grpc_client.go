package rpcdot

import (
	"github.com/scryinfo/dot/dot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConfig struct {
	ServerAddress string
}

type GrpcConnectEx struct {
	conn *grpc.ClientConn
}

func NewGrpcConnectEx(config *GrpcClientConfig, logger *dot.LoggerType) (*GrpcConnectEx, func(), error) {
	conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	d := &GrpcConnectEx{conn: conn}

	return d, func() {
		if d.conn != nil {
			err := d.conn.Close()
			if err != nil {
				logger.Error().Err(err).Send()
			}
		}
	}, nil
}

func (g *GrpcConnectEx) Client() *grpc.ClientConn {
	return g.conn
}
