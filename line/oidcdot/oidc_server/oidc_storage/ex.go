package oidc_storage

type EnumTo[TProto any, TGo any] interface {
	ToGoType(pEnum TProto) TGo
	ToProto(gEnum TGo) TProto
}
