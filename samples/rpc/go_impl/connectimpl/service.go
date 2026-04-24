package connectimpl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/rpcdot"
	apiv1 "github.com/scryinfo/dot/samples/rpc/go_out/connect/api/v1"
	"github.com/scryinfo/dot/samples/rpc/go_out/connect/api/v1/apiv1connect"
)

func NewHiService(mux *rpcdot.ConnectHttpServerMux, logger *dot.LoggerType) *HiService {
	d := &HiService{logger: logger, name: "hi server1"}

	path, handle := apiv1connect.NewHiServiceHandler(d)
	mux.Handle(path, handle)
	return d
}

type HiService struct {
	logger *dot.LoggerType
	name   string
}

// BothStream implements [apiv1connect.HiServiceHandler].
func (p *HiService) BothStream(ctx context.Context, stream *connect.BidiStream[apiv1.BothStreamRequest, apiv1.BothStreamResponse]) error {
	p.logger.Info().Str(p.name, "BothSides").Send()
	count := int64(0)
	for {
		count++
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		req, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			p.logger.Info().Msg("BothStream EOF")
		} else if err != nil {
			p.logger.Error().Err(err).Send()
			return err
		}
		if req != nil {
			res := &apiv1.BothStreamResponse{Greeting: req.Greeting + "  " + strconv.FormatInt(count, 10)}
			err = stream.Send(res)
			if err != nil {
				p.logger.Error().Err(err).Send()
				return err
			}
		} else {
			p.logger.Error().Msg("receive data is nil")
			return nil
		}
	}

}

// ClientStream implements [apiv1connect.HiServiceHandler].
func (p *HiService) ClientStream(ctx context.Context, clientStream *connect.ClientStream[apiv1.ClientStreamRequest]) (*connect.Response[apiv1.ClientStreamResponse], error) {
	p.logger.Info().Str(p.name, "ClientStream").Send()
	count := int64(0)
	messages := strings.Builder{}
	for {
		count++
		select {
		case <-ctx.Done():
			return &connect.Response[apiv1.ClientStreamResponse]{}, nil
		default:
		}
		ok := clientStream.Receive()
		if !ok {
			break
		}
		m := ""
		if clientStream.Msg() != nil {
			m = clientStream.Msg().Greeting
		}
		_, err := fmt.Fprintf(&messages, "%d:%s\r\n", count, m)
		if err != nil {
			return connect.NewResponse(&apiv1.ClientStreamResponse{
				Greeting: err.Error(),
			},
			), err
		}
	}
	res := &apiv1.ClientStreamResponse{Greeting: messages.String()}
	return connect.NewResponse(res), nil
}

// Hi implements [apiv1connect.HiServiceHandler].
func (p *HiService) Hi(_ context.Context, req *connect.Request[apiv1.HiRequest]) (*connect.Response[apiv1.HiResponse], error) {
	p.logger.Info().Str(p.name, req.Msg.Name).Send()
	return connect.NewResponse(&apiv1.HiResponse{Name: p.name}), nil
}

// ServerStream implements [apiv1connect.HiServiceHandler].
func (p *HiService) ServerStream(_ context.Context, req *connect.Request[apiv1.ServerStreamRequest], serverStream *connect.ServerStream[apiv1.ServerStreamResponse]) error {
	p.logger.Info().Str(p.name, "ServerStream").Send()

	res := &apiv1.ServerStreamResponse{Reply: req.Msg.Greeting + " ServerStream"}
	err := serverStream.Send(res) //can multi-send
	if errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

// Write implements [apiv1connect.HiServiceHandler].
func (p *HiService) Write(_ context.Context, req *connect.Request[apiv1.WriteRequest]) (*connect.Response[apiv1.WriteResponse], error) {
	p.logger.Info().Str(p.name, req.Msg.Data).Send()
	res := &apiv1.WriteResponse{Data: "Return : " + req.Msg.Data}
	return connect.NewResponse(res), nil
}
