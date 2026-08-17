// see badger-server
package pebbleservice

import (
	"context"

	"github.com/cockroachdb/pebble/v2"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	kvv1 "github.com/scryinfo/dot/line/db/pebble_service/kv_gen/connect/kv/v1"
	"github.com/scryinfo/dot/line/rpcdot"
	"google.golang.org/grpc"
)

var _ kvv1.KvServiceServer = (*PebbleGrpcService)(nil)

type PebbleGrpcService struct {
	kvv1.UnimplementedKvServiceServer
	db  *pebble2dot.Pebble2
	opt *pebble.WriteOptions
}

func NewPebbleGrpcService(db *pebble2dot.Pebble2, grpcServer *rpcdot.GrpcServer) *PebbleGrpcService {
	service := &PebbleGrpcService{db: db, opt: db.DefaultWriteOpt()}
	grpcServer.ResisterOrLazy(func(server *grpc.Server) {
		server.RegisterService(&kvv1.KvService_ServiceDesc, service)
	})
	return service
}

func (p *PebbleGrpcService) Set(ctx context.Context, req *kvv1.SetRequest) (*kvv1.SetResponse, error) {
	if err := p.db.Db().Set(req.Key, NewValueTTL(req.Value, TtlTime(req.TtlSeconds)).AsBytes(), p.opt); err != nil {
		return nil, err
	}
	return &kvv1.SetResponse{}, nil
}
func (p *PebbleGrpcService) SetSync(ctx context.Context, req *kvv1.SetRequestSync) (*kvv1.SetResponse, error) {
	opt := p.opt
	switch req.Sync {
	case kvv1.SyncMode_SYNC:
		opt = pebble.Sync
	case kvv1.SyncMode_ASYNC:
		opt = pebble.NoSync
	}
	if err := p.db.Db().Set(req.Key, NewValueTTL(req.Value, TtlTime(req.TtlSeconds)).AsBytes(), opt); err != nil {
		return nil, err
	}
	return &kvv1.SetResponse{}, nil
}
func (p *PebbleGrpcService) SetKvPair(ctx context.Context, req *kvv1.KvPair) (*kvv1.SetResponse, error) {
	if err := p.db.Db().Set(req.Key, NewValueKvPair(req.Value, req.ExpireAt).AsBytes(), p.opt); err != nil {
		return nil, err
	}
	return &kvv1.SetResponse{}, nil
}
func (p *PebbleGrpcService) SetKvPairSync(ctx context.Context, req *kvv1.KvPairSync) (*kvv1.SetResponse, error) {
	opt := p.opt
	switch req.Sync {
	case kvv1.SyncMode_SYNC:
		opt = pebble.Sync
	case kvv1.SyncMode_ASYNC:
		opt = pebble.NoSync
	}
	if err := p.db.Db().Set(req.Key, NewValueKvPair(req.Value, req.ExpireAt).AsBytes(), opt); err != nil {
		return nil, err
	}
	return &kvv1.SetResponse{}, nil
}
func (p *PebbleGrpcService) Get(ctx context.Context, req *kvv1.GetRequest) (*kvv1.GetResponse, error) {
	value, closer, err := p.db.Db().Get(req.Key)
	if err == pebble.ErrNotFound {
		return &kvv1.GetResponse{Value: nil, Found: false}, nil
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
		return &kvv1.GetResponse{Value: nil, Found: false}, nil
	} else {
		return &kvv1.GetResponse{Value: kvValue.Value(), Found: true}, nil
	}
}

func (p *PebbleGrpcService) Delete(ctx context.Context, req *kvv1.DeleteRequest) (*kvv1.DeleteResponse, error) {
	if err := p.db.Db().Delete(req.Key, p.opt); err != nil {
		return nil, err
	}
	return &kvv1.DeleteResponse{}, nil
}
func (p *PebbleGrpcService) DeleteSync(ctx context.Context, req *kvv1.DeleteRequestSync) (*kvv1.DeleteResponse, error) {
	opt := p.opt
	switch req.Sync {
	case kvv1.SyncMode_SYNC:
		opt = pebble.Sync
	case kvv1.SyncMode_ASYNC:
		opt = pebble.NoSync
	}
	if err := p.db.Db().Delete(req.Key, opt); err != nil {
		return nil, err
	}
	return &kvv1.DeleteResponse{}, nil
}
func (p *PebbleGrpcService) Scan(ctx context.Context, req *kvv1.ScanRequest) (*kvv1.ScanResponse, error) {

	iterOpts := &pebble.IterOptions{
		LowerBound: req.Start,
		UpperBound: req.End,
	}
	iter, err := p.db.Db().NewIter(iterOpts)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	count := int32(0)
	list := make([]*kvv1.KvPair, 0, req.Limit)
	for iter.First(); iter.Valid(); iter.Next() {
		count++
		if count >= req.Limit {
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
	return &kvv1.ScanResponse{List: list}, nil
}

func (p *PebbleGrpcService) BatchSet(ctx context.Context, req *kvv1.BatchSetRequest) (*kvv1.BatchSetResponse, error) {
	tx := p.db.Db().NewBatch()
	defer tx.Close()
	for _, e := range req.Entries {
		kvValue := NewValueTTL(e.Value, TtlTime(e.TtlSeconds))
		err := tx.Set(e.Key, kvValue.AsBytes(), p.opt)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(p.opt); err != nil {
		return nil, err
	}
	return &kvv1.BatchSetResponse{}, nil
}

func (p *PebbleGrpcService) BatchSetSync(ctx context.Context, req *kvv1.BatchSetRequestSync) (*kvv1.BatchSetResponse, error) {
	tx := p.db.Db().NewBatch()
	defer tx.Close()
	opt := p.opt
	switch req.Sync {
	case kvv1.SyncMode_SYNC:
		opt = pebble.Sync
	case kvv1.SyncMode_ASYNC:
		opt = pebble.NoSync
	}
	for _, e := range req.Entries {
		kvValue := NewValueTTL(e.Value, TtlTime(e.TtlSeconds))
		err := tx.Set(e.Key, kvValue.AsBytes(), opt)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(opt); err != nil {
		return nil, err
	}
	return &kvv1.BatchSetResponse{}, nil
}
func (p *PebbleGrpcService) BatchSetKvPair(ctx context.Context, req *kvv1.BatchSetKvPairRequest) (*kvv1.BatchSetKvPairResponse, error) {
	tx := p.db.Db().NewBatch()
	defer tx.Close()
	for _, e := range req.Entries {
		kvValue := NewValueKvPair(e.Value, e.ExpireAt)
		err := tx.Set(e.Key, kvValue.AsBytes(), p.opt)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(p.opt); err != nil {
		return nil, err
	}
	return &kvv1.BatchSetKvPairResponse{}, nil
}

func (p *PebbleGrpcService) BatchSetKvPairSync(ctx context.Context, req *kvv1.BatchSetKvPairRequestSync) (*kvv1.BatchSetKvPairResponse, error) {
	tx := p.db.Db().NewBatch()
	defer tx.Close()
	opt := p.opt
	switch req.Sync {
	case kvv1.SyncMode_SYNC:
		opt = pebble.Sync
	case kvv1.SyncMode_ASYNC:
		opt = pebble.NoSync
	}
	for _, e := range req.Entries {
		kvValue := NewValueKvPair(e.Value, e.ExpireAt)
		err := tx.Set(e.Key, kvValue.AsBytes(), opt)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(opt); err != nil {
		return nil, err
	}
	return &kvv1.BatchSetKvPairResponse{}, nil
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
