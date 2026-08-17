package kits

import (
	"time"
	"unsafe"
)

type _Tss struct{}

var Tss = _Tss{}

// seconds
func (c _Tss) TsSeconds() int64 {
	return time.Now().Unix()
}

// microseconds
func (c _Tss) TsMicroseconds() int64 {
	return time.Now().UnixMicro()
}

// nanoseconds
func (c _Tss) TsNanoseconds() int64 {
	return time.Now().UnixNano()
}

// milliseconds
func (c _Tss) TsMilliseconds() int64 {
	return time.Now().UnixMilli()
}

// inline
func StringToBytes(data string) []byte {
	return unsafe.Slice(unsafe.StringData(data), len(data))
}

// inline
func BytesToString(data []byte) string {
	return unsafe.String(unsafe.SliceData(data), len(data))
}
