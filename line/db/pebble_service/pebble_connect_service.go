// see badger-server
package pebbleservice

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/pebble/v2"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	kvv1 "github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1"
	"github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1/kvv1connect"
	"github.com/scryinfo/dot/line/rpcdot"
)

var _ kvv1connect.KvServiceHandler = (*PebbleConnectService)(nil)

type PebbleConnectService struct {
	// kvv1.UnimplementedKVServerServer
	db *pebble.DB
}

// func NewHiService(mux *rpcdot.ConnectHttpServerMux, logger *dot.LoggerType, conf *HiServiceConfig) *HiService {
// 	d := &HiService{logger: logger, name: conf.Name}

// 	path, handle := apiv1connect.NewHiServiceHandler(d)
// 	mux.Handle(path, handle)
// 	return d
// }

func NewPebbleConnectService(db *pebble2dot.Pebble2, mux *rpcdot.ConnectHttpServerMux) *PebbleConnectService {
	service := &PebbleConnectService{db: db.Db()}
	path, handle := kvv1connect.NewKvServiceHandler(service)
	mux.Handle(path, handle)
	return service
}

func (p *PebbleConnectService) Set(ctx context.Context, req *connect.Request[kvv1.SetRequest]) (*connect.Response[kvv1.SetResponse], error) {
	if err := p.db.Set(req.Msg.Key, NewValueTTL(req.Msg.Value, TtlTime(req.Msg.TtlSeconds)).AsBytes(), pebble.Sync); err != nil {
		return nil, err
	}
	return &connect.Response[kvv1.SetResponse]{}, nil
}
func (p *PebbleConnectService) SetKvPair(ctx context.Context, req *connect.Request[kvv1.KvPair]) (*connect.Response[kvv1.SetResponse], error) {
	if err := p.db.Set(req.Msg.Key, NewValueKvPair(req.Msg.Value, req.Msg.ExpireAt).AsBytes(), pebble.Sync); err != nil {
		return nil, err
	}
	return &connect.Response[kvv1.SetResponse]{}, nil
}
func (p *PebbleConnectService) Get(ctx context.Context, req *connect.Request[kvv1.GetRequest]) (*connect.Response[kvv1.GetResponse], error) {
	value, closer, err := p.db.Get(req.Msg.Key)
	if err == pebble.ErrNotFound {
		return &connect.Response[kvv1.GetResponse]{Msg: &kvv1.GetResponse{Value: nil, Found: false}}, nil
	} else if err != nil {
		return nil, err
	}
	var copyValue []byte
	copyValue = append(copyValue, value...)
	if err = closer.Close(); err != nil {
		return nil, err
	}

	kvValue := KvValue(copyValue)
	if kvValue.HasExpire() {
		return &connect.Response[kvv1.GetResponse]{Msg: &kvv1.GetResponse{Value: nil, Found: false}}, nil
	} else {
		return &connect.Response[kvv1.GetResponse]{Msg: &kvv1.GetResponse{Value: kvValue.Value(), Found: true}}, nil
	}
}

func (p *PebbleConnectService) Delete(ctx context.Context, req *connect.Request[kvv1.DeleteRequest]) (*connect.Response[kvv1.DeleteResponse], error) {
	if err := p.db.Delete(req.Msg.Key, pebble.Sync); err != nil {
		return nil, err
	}
	return &connect.Response[kvv1.DeleteResponse]{}, nil
}

func (p *PebbleConnectService) Scan(ctx context.Context, req *connect.Request[kvv1.ScanRequest]) (*connect.Response[kvv1.ScanResponse], error) {

	iterOpts := &pebble.IterOptions{
		LowerBound: req.Msg.Start,
		UpperBound: req.Msg.End,
	}
	iter, err := p.db.NewIter(iterOpts)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	count := int32(0)
	list := make([]*kvv1.KvPair, 0, req.Msg.Limit)
	for iter.First(); iter.Valid(); iter.Next() {
		count++
		if count >= req.Msg.Limit {
			break
		}
		kvValue := KvValue(iter.Value())
		if kvValue.HasExpire() {
			count--
			continue
		}
		v := kvv1.KvPair{
			ExpireAt: kvValue.ExpireAt(),
		}
		v.Key = append(v.Key, iter.Key()...)
		v.Value = append(v.Value, kvValue.Value()...)
		list = append(list, &v)
	}
	return &connect.Response[kvv1.ScanResponse]{Msg: &kvv1.ScanResponse{List: list}}, nil
}

func (p *PebbleConnectService) BatchSet(ctx context.Context, req *connect.Request[kvv1.BatchSetRequest]) (*connect.Response[kvv1.BatchSetResponse], error) {
	tx := p.db.NewBatch()
	defer tx.Close()
	for _, e := range req.Msg.Entries {
		kvValue := NewValueTTL(e.Value, TtlTime(e.TtlSeconds))
		err := tx.Set(e.Key, kvValue.AsBytes(), pebble.Sync)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(pebble.Sync); err != nil {
		return nil, err
	}
	return &connect.Response[kvv1.BatchSetResponse]{}, nil
}
func (p *PebbleConnectService) BatchSetKvPair(ctx context.Context, req *connect.Request[kvv1.BatchSetKvPairRequest]) (*connect.Response[kvv1.BatchSetKvPairResponse], error) {
	tx := p.db.NewBatch()
	defer tx.Close()
	for _, e := range req.Msg.Entries {
		kvValue := NewValueKvPair(e.Value, e.ExpireAt)
		err := tx.Set(e.Key, kvValue.AsBytes(), pebble.Sync)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(pebble.Sync); err != nil {
		return nil, err
	}
	return &connect.Response[kvv1.BatchSetKvPairResponse]{}, nil
}

// func main() {
// 	flag.Parse()
// 	opt := badger.DefaultOptions(*dataDir)
// 	db, err := badger.Open(opt)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	go gcLoop(db)

// 	lis, err := net.Listen("tcp", *grpcAddr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	srv := grpc.NewServer()
// 	kvv1.RegisterKVServerServer(srv, newKVServer(db))

// 	go func() {
// 		fmt.Printf("gRPC listen: %s\n", *grpcAddr)
// 		_ = srv.Serve(lis)
// 	}()

// 	ctx := context.Background()
// 	mux := runtime.NewServeMux()
// 	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
// 	if err := kvv1.RegisterKVServerHandlerFromEndpoint(ctx, mux, *grpcAddr, dialOpts); err != nil {
// 		panic(err)
// 	}

// 	go func() {
// 		fmt.Printf("HTTP gateway listen: %s\n", *httpAddr)
// 		_ = http.ListenAndServe(*httpAddr, mux)
// 	}()

// 	sigCh := make(chan os.Signal, 1)
// 	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
// 	<-sigCh
// 	fmt.Println("Shutting down...")
// 	srv.GracefulStop()
// }
