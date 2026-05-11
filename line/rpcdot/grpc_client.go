package rpcdot

import (
	"github.com/scryinfo/dot/dot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConfig struct {
	ServerAddress string
}

type GrpcClientEx struct {
	conn *grpc.ClientConn
}

func NewGrpcClientEx(config *GrpcClientConfig, logger *dot.LoggerType) (*GrpcClientEx, func(), error) {
	conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	d := &GrpcClientEx{conn: conn}

	return d, func() {
		if d.conn != nil {
			err := d.conn.Close()
			if err != nil {
				logger.Error().Err(err).Send()
			}
		}
	}, nil
}

func (g *GrpcClientEx) Client() *grpc.ClientConn {
	return g.conn
}
