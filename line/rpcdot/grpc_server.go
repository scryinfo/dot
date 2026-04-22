package rpcdot

import (
	"google.golang.org/grpc"
)

func NewGrpcServer() *grpc.Server {
	return grpc.NewServer()
}
