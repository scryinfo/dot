package pebbleservice

import (
	"encoding/binary"
	"time"
)

type KvValue []byte

type KvKey []byte

type KvEntry struct {
	Key   KvKey
	Value KvValue
}

// TtlTime is the number of seconds until the entry expires.
type TtlTime uint32

const (
	typeOffset  = 1
	valueOffset = 9
)

func NewValueTTL(value []byte, ttlSeconds TtlTime) KvValue {
	v := make([]byte, 1, len(value)+valueOffset)
	expireAt := uint64(0)
	if ttlSeconds > 0 {
		v[0] = byte(1)
		now := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		expireAt = uint64(now.Unix())
	} else {
		v[0] = byte(0)
	}
	v = binary.LittleEndian.AppendUint64(v, expireAt)
	v = append(v, value...)
	return KvValue(v)
}

func NewValueKvPair(value []byte, expireAt uint64) KvValue {
	v := make([]byte, 1, len(value)+valueOffset)
	if expireAt > 0 {
		v[0] = byte(1)
	} else {
		v[0] = byte(0)
	}
	v = binary.LittleEndian.AppendUint64(v, expireAt)
	v = append(v, value...)
	return KvValue(v)
}

func (v KvValue) Value() []byte {
	return v[valueOffset:]
}

func (v KvValue) ExpireAt() uint64 {
	return binary.LittleEndian.Uint64(v[typeOffset:])
}

func (v KvValue) HasExpire() bool {
	return v[0] > 0 && int64(binary.LittleEndian.Uint64(v[typeOffset:])) < time.Now().Unix()
}

func (v KvValue) HasTtl() bool {
	return v[0] == 1
}

func (v KvValue) AsBytes() []byte {
	return v
}

func (k KvKey) AsBytes() []byte {
	return k
}
