//go:generate goverter gen ./...
package db

import (
	datav1 "github.com/scryinfo/dot/samples/db/sqlcs/proto/data_gen/connect/data/v1"
)

// goverter:converter
// goverter:output:file ./gen_to/gen_sqlc_proto.go
// goverter:ignoreUnexported
// goverter:ignoreMissing
// goverter:extend PgTextToString
// goverter:extend StringtToPgText
type UserConverter interface {
	ToProto(src *User) *datav1.User
	ToProtoSlice(src []*User) []*datav1.User
	ToSqlc(src *datav1.User) *User
	ToSqlcSlice(src []*datav1.User) []*User
}
