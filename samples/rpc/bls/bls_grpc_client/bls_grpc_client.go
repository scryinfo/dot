// Scry Info.  All rights reserved.
// license that can be found in the license file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/etcddot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	apiv1grpc "github.com/scryinfo/dot/samples/rpc/go_out/gogrpc/api/v1"
	"github.com/scryinfo/scryg/sutils/ssignal"
)

type Line struct {
	SConfig         *sconfig.SConfig
	Logger          *dot.LoggerType
	HiServiceClient apiv1grpc.HiServiceClient
}

type LineConfig struct {
	Log            dot.LogConfig               `json:"log" toml:"log" yaml:"log" mapstructure:"log"`
	GrpcClientEtcd rpcdot.GrpcClientEtcdConfig `json:"grpc_client_etcd" toml:"grpc_client_etcd" yaml:"grpc_client_etcd" mapstructure:"grpc_client_etcd"`
	EtcdClient     etcddot.ClientConfig        `json:"etcd_client" toml:"etcd_client" yaml:"etcd_client" mapstructure:"etcd_client"`
}

func NewLineConfig(config *sconfig.SConfig) (*LineConfig, error) {
	lineConfig, err := sconfig.NewLineConfig[LineConfig](config)
	if err != nil {
		return nil, err
	}
	return sconfig.GenerateConfigWithArgs(config, lineConfig)
}

var LineSet = wire.NewSet(
	wire.Struct(new(Line), "*"),
	wire.FieldsOf(new(*LineConfig), "Log", "GrpcClientEtcd", "EtcdClient"),
	NewLineConfig,
	line.SconfigNewConfig,
	wire.Bind(new(dot.SConfig), new(*sconfig.SConfig)),
	line.CertificateNewBaseCertificate,
	dot.NewLogger,
	line.ContextexNewContextEx,
	line.RpcdotNewGrpcClientEtcd,
	NewHiServiceClient,
	line.EtcddotNewClient,
)

func NewHiServiceClient(clientEtcd *rpcdot.GrpcClientEtcd) apiv1grpc.HiServiceClient {
	return apiv1grpc.NewHiServiceClient(clientEtcd.Client())
}

func main() {
	// dot.InitLogger(new(dot.TestLogConfig()))
	line, clean, err := InitializeService()
	if err != nil {
		if line != nil {
			dot.Logger.Error().Err(err).Msg("initialize service failed")
		} else {
			dot.Logger.Info().Msg(err.Error())
			fmt.Printf("\n\n")
		}
		return
	}
	if clean != nil {
		defer clean()
	}

	line.Logger.Info().Msg("dot ok")
	//second step ....
	res, err := line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "ttt"})
	if err != nil {
		line.Logger.Error().Err(err).Msg("hi failed")
		return
	}
	line.Logger.Info().Msg(res.Name)
	res, err = line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "ttt2"})
	if err != nil {
		line.Logger.Error().Err(err).Msg("hi failed")
		return
	}
	line.Logger.Info().Msg(res.Name)
	res, err = line.HiServiceClient.Hi(context.Background(), &apiv1grpc.HiRequest{Name: "ttt3"})
	if err != nil {
		line.Logger.Error().Err(err).Msg("hi failed")
		return
	}
	line.Logger.Info().Msg(res.Name)

	ssignal.WaitCtrlC(func(s os.Signal) bool { //third wait for exit
		return false
	})
}
