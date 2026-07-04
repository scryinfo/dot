package rocksdbbench

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"connectrpc.com/connect"
	"github.com/cockroachdb/pebble/v2"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	pebbleservice "github.com/scryinfo/dot/line/db/pebble_service"
	kvv1 "github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
)

func newLogger() *dot.LoggerType {
	conf := dot.TestLogConfig()
	conf.AddStdout = false
	conf.SetSlog = false
	conf.Level = "error"
	return dot.NewLogger(&conf)
}

func BenchmarkConnect(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	logger.Info().Msgf("rocksdb path: %s", filepath.Join(sourcePath, "temp/gorocksdb"))

	baseCertificate := line.CertificateNewBaseCertificate(logger)
	sconfig := sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath)
	clientEx, err := rpcdot.NewHttpClientEx(&rpcdot.HttpClientConfig{
		ForceAttemptHTTP2:   true,
		DisableCompression:  false,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		ServerAddress:       "http://localhost:8090",
		Tls: rpcdot.TlsConfig{
			Mode:     "none",
			Cert:     "",
			Key:      "",
			RootCert: "",
			PeerCert: "",
		},
	}, sconfig, baseCertificate, logger)
	if err != nil {
		b.Fatal(err)
	}
	pebbleConnect := pebbleservice.NewPebbleConnectClient(clientEx)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pebbleConnect.Set(context.Background(), &connect.Request[kvv1.SetRequest]{
			Msg: &kvv1.SetRequest{
				Key:   binary.LittleEndian.AppendUint32(nil, uint32(i)),
				Value: []byte("value"),
			},
		})
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}

func BenchmarkGrpc(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()

	baseCertificate := line.CertificateNewBaseCertificate(logger)
	sconfig := sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath)

	clientEx, clearer, err := rpcdot.NewGrpcClientEx(&rpcdot.GrpcClientConfig{
		ServerAddress: "localhost:8091",
		Tls: rpcdot.TlsConfig{
			Mode:     "none",
			Cert:     "",
			Key:      "",
			RootCert: "",
			PeerCert: "",
		},
	}, sconfig, logger, baseCertificate)
	if err != nil {
		b.Fatal(err)
	}
	if clearer != nil {
		defer clearer()
	}
	pebbleConnect := pebbleservice.NewPebbleGrpcClient(clientEx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pebbleConnect.Set(context.Background(), &kvv1.SetRequest{
			Key:   binary.LittleEndian.AppendUint32([]byte{}, uint32(i)),
			Value: []byte("value"),
		})
		if err != nil {
			b.Fatal(err)
		}

	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}

// todo
func BenchmarkGrpcClientConnect(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	baseCertificate := line.CertificateNewBaseCertificate(logger)
	sconfig := sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath)

	clientEx, clearer, err := rpcdot.NewGrpcClientEx(&rpcdot.GrpcClientConfig{
		ServerAddress: "localhost:8092",
		Tls: rpcdot.TlsConfig{
			Mode:     "none",
			Cert:     "",
			Key:      "",
			RootCert: "",
			PeerCert: "",
		},
	}, sconfig, logger, baseCertificate)
	if err != nil {
		b.Fatal(err)
	}
	if clearer != nil {
		defer clearer()
	}
	pebbleConnect := pebbleservice.NewPebbleGrpcClient(clientEx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pebbleConnect.Set(context.Background(), &kvv1.SetRequest{
			Key:   binary.LittleEndian.AppendUint32([]byte{}, uint32(i)),
			Value: []byte("value"),
		})
		if err != nil {
			b.Fatal(err)
		}

	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}

func BenchmarkPebbleEmbedded(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	config := pebble2dot.Pebble2Config{
		DbPath: filepath.Join(sourcePath, "temp/pebble"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		if err != nil {
			b.Fatal(err)
		}
	}
	db, cleaner, err := pebble2dot.NewPebble2(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
	if err != nil {
		b.Fatal(err)
	}
	defer cleaner()

	b.ResetTimer()
	opt := &pebble.WriteOptions{}

	for i := 0; i < b.N; i++ {
		err := db.Db().Set(binary.LittleEndian.AppendUint32(nil, uint32(i)), []byte("value"), opt)
		if err != nil {
			b.Fatal(err)
		}

	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}
