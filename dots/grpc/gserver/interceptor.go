// Scry Info.  All rights reserved.
// license that can be found in the license file.

package gserver

import (
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = status.Errorf(codes.Internal, "Panic err: %v", e)
				dot.Logger().Errorln("", zap.Error(err))
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = status.Errorf(codes.Internal, "Panic err: %v", e)
				dot.Logger().Errorln("", zap.Error(err))
			}
		}()
		return handler(srv, stream)
	}
}
