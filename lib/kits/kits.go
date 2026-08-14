package kits

import "time"

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
