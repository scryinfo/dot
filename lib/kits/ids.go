package kits

import (
	"crypto/rand"
	"encoding/base64"
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
func (c _Ids) NewAuthCode() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic("rand.Read: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

// inline
func UnsafeToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}
