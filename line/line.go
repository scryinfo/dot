package line

import (
	"github.com/scryinfo/dot/line/certificate"
	contextex "github.com/scryinfo/dot/line/context_ex"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/dot/line/db/badgerdot/backup"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/db/redis_client"
	"github.com/scryinfo/dot/line/etcddot"
	"github.com/scryinfo/dot/line/gindot"
	"github.com/scryinfo/dot/line/rpcdot"
	"github.com/scryinfo/dot/line/sconfig"
)

var (
	CertificateNewBaseCertificate = certificate.NewBaseCertificate
	CertificateNewEcdsa           = certificate.NewEcdsa
	CertificateNewEd25519         = certificate.NewEd25519
	CertificateNewRsa             = certificate.NewRsa
	CertificateNewSm2             = certificate.NewSm2

	ContextexNewContextEx       = contextex.NewContextEx
	DbBadgerdotNewBadgerDot     = badgerdot.NewBadgerDot
	DbBadgerdotNewBackup        = backup.NewDbBackup
	DbPebble2dotNewPebble2      = pebble2dot.NewPebble2
	DbRedisClientNewRedisClient = redis_client.NewRedisClient
	// dont difine rocksdb
	EtcddotNewClient = etcddot.NewClient
	EtcddotNewServer = etcddot.NewServer

	GindotNewGinDot = gindot.NewGinDot
	GindotNewRouter = gindot.NewRouter
	GindotNewUi     = gindot.NewUi

	RpcdotNewBothHttpServer       = rpcdot.NewBothHttpServer
	RpcdotNewHttpClientEx         = rpcdot.NewHttpClientEx
	RpcdotNewHttpClientEtcd       = rpcdot.NewHttpClientEtcd
	RpcdotNewHandlerMiddle        = rpcdot.NewHandlerMiddle
	RpcdotNewConnectHttpServerMux = rpcdot.NewConnectHttpServerMux
	RpcdotNewConnetServer         = rpcdot.NewConnetServer
	RpcdotNewConnectServerEtcd    = rpcdot.NewConnectServerEtcd
	RpcdotNewGrpcClientEtcd       = rpcdot.NewGrpcClientEtcd
	RpcdotNewGrpcClientEx         = rpcdot.NewGrpcClientEx
	RpcdotNewGrpcServer           = rpcdot.NewGrpcServer

	SconfigNewConfig      = sconfig.NewConfig
	SconfigNewDataWithNet = sconfig.NewDataWithNet
)
