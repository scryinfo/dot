// Scry Info.  All rights reserved.
// license that can be found in the license file.

package nobl

import (
	"context"
	"io"
	"strconv"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/sample/grpc/go_out/hidot"
	"go.uber.org/zap"
)

var _ hidot.HiDotServer = (*HiServer)(nil)

const (
	HiServerTypeID = "hiserver"
)

type config struct {
	Name string `json:"name"`
}

type HiServer struct {
	hidot.UnimplementedHiDotServer
	ServerNobl gserver.ServerNobl `dot:""`
	conf       config

	ctx       context.Context
	cancelFun context.CancelFunc
}

func newHiServer(conf interface{}) (dot.Dot, error) {
	var err error = nil
	var bs []byte = nil
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &config{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d := &HiServer{
		conf: *dconf,
	}
	d.ctx, d.cancelFun = context.WithCancel(context.Background())

	return d, err
}

func (c *HiServer) Hi(ctx context.Context, req *hidot.HiReq) (*hidot.HiRes, error) {
	dot.Logger().Infoln("HiServer", zap.String(c.conf.Name, req.Name))
	res := &hidot.HiRes{Name: c.conf.Name}
	return res, nil
}

func (c *HiServer) Write(ctx context.Context, req *hidot.WriteReq) (*hidot.WriteRes, error) {
	dot.Logger().Infoln("HiServer", zap.String(c.conf.Name, req.Data))
	res := &hidot.WriteRes{Data: "Return : " + req.Data}
	return res, nil
}

func (c *HiServer) ServerStream(req *hidot.HelloRequest, serverStream hidot.HiDot_ServerStreamServer) error {
	dot.Logger().Infoln("HiServer", zap.String(c.conf.Name, "ServerStream"))

	res := &hidot.HelloResponse{Reply: req.Greeting + " ServerStream"}
	err := serverStream.Send(res)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	//req.Reset()
	//err = serverStream.RecvMsg(req)
	//if err == io.EOF {
	//	return nil
	//} else if err != nil {
	//	return err
	//}
	return nil
}

func (c *HiServer) ClientStream(clientStream hidot.HiDot_ClientStreamServer) error {
	dot.Logger().Infoln("HiServer", zap.String(c.conf.Name, "ClientStream"))
	count := int64(0)
	for {
		count++
		ctx := clientStream.Context()
		select {
		case <-ctx.Done():
			return nil
		case <-c.ctx.Done():
			return nil
		default:

		}
		req, err := clientStream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		res := &hidot.HelloResponse{Reply: req.Greeting + "  " + strconv.FormatInt(count, 10)}
		clientStream.SendMsg(res)
	}
}

func (c *HiServer) BothSides(server hidot.HiDot_BothSidesServer) error {
	dot.Logger().Infoln("HiServer", zap.String(c.conf.Name, "BothSides"))
	count := int64(0)
	for {
		count++
		ctx := server.Context()
		select {
		case <-ctx.Done():
			return nil
		case <-c.ctx.Done():
			return nil
		default:

		}
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		res := &hidot.HelloResponse{Reply: req.Greeting + "  " + strconv.FormatInt(count, 10)}
		server.SendMsg(res)
	}
}

func (c *HiServer) Start(ignore bool) error {
	hidot.RegisterHiDotServer(c.ServerNobl.Server(), c)
	return nil
}

// HiServerTypeLives make all type lives
func HiServerTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{{
		Meta: dot.Metadata{TypeID: HiServerTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newHiServer(conf)
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveID:    HiServerTypeID,
				RelyLives: map[string]dot.LiveID{"ServerNobl": gserver.ServerNoblTypeID},
			},
		},
	}}
	lives = append(lives, gserver.ServerNoblTypeLives()...)
	return lives
}
