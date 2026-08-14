package kits

import (
	"time"
	"unsafe"

	"github.com/rs/xid"
)

type _Ids struct{}

var Ids = _Ids{}

// seconds
func (c _Ids) Ts() int64 {
	return time.Now().Unix()
}

// inline
func (c _Ids) NewXId() string {
	return xid.New().String()
}

// inline
func UnsafeToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}
