package rocksdbbench

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"sync/atomic"
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

// one test result
// make bench
// "go: /home/peace/bin/go/bin/go"
// /home/peace/bin/go/bin/go test -bench=. -run=^$ -count=3 -benchmem -benchtime=3s
// goos: linux
// goarch: amd64
// pkg: github.com/scryinfo/dot/samples/db/pebble_server_bench
// cpu: AMD Ryzen 9 7900X 12-Core Processor
// BenchmarkConnect-24              	  217110	     16799 ns/op	     59527 tps	   17039 B/op	     119 allocs/op
// BenchmarkConnect-24              	  225261	     16467 ns/op	     60727 tps	   16584 B/op	     119 allocs/op
// BenchmarkConnect-24              	  222133	     16335 ns/op	     61217 tps	   16196 B/op	     119 allocs/op
// BenchmarkGrpc-24                 	  346990	     10067 ns/op	     99345 tps	    5271 B/op	      80 allocs/op
// BenchmarkGrpc-24                 	  362470	     10009 ns/op	     99912 tps	    5271 B/op	      80 allocs/op
// BenchmarkGrpc-24                 	  390178	      9891 ns/op	    101107 tps	    5267 B/op	      80 allocs/op
// BenchmarkGrpcClientConnect-24    	   94484	     37726 ns/op	     26508 tps	    8007 B/op	      94 allocs/op
// BenchmarkGrpcClientConnect-24    	   95007	     38243 ns/op	     26149 tps	    7890 B/op	      93 allocs/op
// BenchmarkGrpcClientConnect-24    	   95888	     38075 ns/op	     26265 tps	    7905 B/op	      93 allocs/op
// BenchmarkPebbleEmbedded-24       	 5532043	      1308 ns/op	    824999 tps	      14 B/op	       1 allocs/op
// BenchmarkPebbleEmbedded-24       	 5432888	      1328 ns/op	    818955 tps	      14 B/op	       1 allocs/op
// BenchmarkPebbleEmbedded-24       	 5565430	      1277 ns/op	    845033 tps	      14 B/op	       1 allocs/op
// PASS
// ok  	github.com/scryinfo/dot/samples/db/pebble_server_bench	63.322s

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
	ctx := context.Background()
	var requestCounter uint32 = 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint32(&requestCounter, 1)
			_, err := pebbleConnect.Set(ctx, &connect.Request[kvv1.SetRequest]{
				Msg: &kvv1.SetRequest{
					Key:   binary.LittleEndian.AppendUint32(nil, id),
					Value: []byte("value"),
				},
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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
	ctx := context.Background()
	var requestCounter uint32 = 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint32(&requestCounter, 1)
			_, err := pebbleConnect.Set(ctx, &kvv1.SetRequest{
				Key:   binary.LittleEndian.AppendUint32([]byte{}, id),
				Value: []byte("value"),
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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
	ctx := context.Background()
	var requestCounter uint32 = 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint32(&requestCounter, 1)
			_, err := pebbleConnect.Set(ctx, &kvv1.SetRequest{
				Key:   binary.LittleEndian.AppendUint32([]byte{}, id),
				Value: []byte("value"),
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}

func BenchmarkPebbleEmbedded(b *testing.B) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	logger := newLogger()
	config := pebble2dot.Pebble2Config{
		DbPath: filepath.Join(sourcePath, "data/pebble_embedded"),
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
	opt := &pebble.WriteOptions{
		Sync: false,
	}
	// var requestCounter uint32 = 0

	for i := 0; i < b.N; i++ {
		// id := atomic.AddUint32(&requestCounter, 1)
		err := db.Db().Set(binary.LittleEndian.AppendUint32(nil, uint32(i)), []byte("value"), opt)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tps")
}
