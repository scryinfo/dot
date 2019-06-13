package nobl

import (
	"context"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/gserver"
	"github.com/scryinfo/dot/sample/grpc/go_out/hidot"
	"go.uber.org/zap"
)

const (
	HiServerTypeId = "hiserver"
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

func (serv *HiServer) Hi(ctx context.Context, req *hidot.ReqData) (*hidot.ResData, error) {
	dot.Logger().Infoln("HiServer", zap.String(serv.conf.Name, req.Name))
	res := &hidot.ResData{Name: serv.conf.Name}
	return res, nil
}

func (serv *HiServer) Start(ignore bool) error {
	hidot.RegisterHiDotServer(serv.ServerNobl.Server(), serv)
	return nil
}

//HiServerTypeLives make all type lives
func HiServerTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: HiServerTypeId, NewDoter: func(conf interface{}) (dot.Dot, error) {
			return newHiServer(conf)
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    HiServerTypeId,
				RelyLives: map[string]dot.LiveId{"ServerNobl": gserver.ServerNoblTypeId},
			},
		},
	}

	lives := []*dot.TypeLives{gserver.ConnsTypeLives(), tl}

	return lives
}
