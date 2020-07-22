// Scry Info.  All rights reserved.
// license that can be found in the license file.

package nobl

import (
	"context"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/sample/grpc/go_out/hidot"
	"go.uber.org/zap"
)

const (
	HiServerTypeID = "hiserver"
)

type config struct {
	Name string `json:"name"`
}

type HiServer struct {
	ServerNobl gserver.ServerNobl `dot:""`
	conf       config
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

	return d, err
}

func (serv *HiServer) Hi(ctx context.Context, req *hidot.HiReq) (*hidot.HiRes, error) {
	dot.Logger().Infoln("HiServer", zap.String(serv.conf.Name, req.Name))
	res := &hidot.HiRes{Name: serv.conf.Name}
	return res, nil
}

func (serv *HiServer) Write(ctx context.Context, req *hidot.WriteReq) (*hidot.WriteRes, error) {
	dot.Logger().Infoln("HiServer", zap.String(serv.conf.Name, req.Data))
	res := &hidot.WriteRes{Data: "Return : " + req.Data}
	return res, nil
}

func (serv *HiServer) Start(ignore bool) error {
	hidot.RegisterHiDotServer(serv.ServerNobl.Server(), serv)
	return nil
}

//HiServerTypeLives make all type lives
func HiServerTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: HiServerTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return newHiServer(conf)
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveID:    HiServerTypeID,
				RelyLives: map[string]dot.LiveID{"ServerNobl": gserver.ServerNoblTypeID},
			},
		},
	}

	lives := []*dot.TypeLives{gserver.ServerNoblTypeLive(), tl}

	return lives
}
