package oidc_impl

import (
	"encoding/binary"
	"time"

	"github.com/google/uuid"

	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
)

// ts microseconds
func NewTs() int64 {
	return time.Now().UnixMicro()
}

// // UUIDToUint64 将 uuid.UUID 拆分为 high 和 low 两个 uint64
// func UUIDToUint64(id uuid.UUID) (high, low uint64) {
// 	// UUID 为 16 字节：前 8 字节为 high (包含 v7 的时间戳)，后 8 字节为 low
// 	high = binary.BigEndian.Uint64(id[0:8])
// 	low = binary.BigEndian.Uint64(id[8:16])
// 	return high, low
// }

// // Uint64ToUUID 将 high 和 low 两个 uint64 还原为 uuid.UUID
// func Uint64ToUUID(high, low uint64) uuid.UUID {
// 	var id uuid.UUID
// 	binary.BigEndian.PutUint64(id[0:8], high)
// 	binary.BigEndian.PutUint64(id[8:16], low)
// 	return id
// }

func NewUuIdV7(req *oidcapiv1.Reqbase) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	req.ReqIdLow, req.ReqIdHigh = binary.BigEndian.Uint64(id[0:8]), binary.BigEndian.Uint64(id[8:16])
	return nil
}

func NewUuidV7Resbase(req *oidcapiv1.Reqbase, res *oidcapiv1.Resbase) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	res.ReqIdLow = req.ReqIdLow
	res.ReqIdHigh = req.ReqIdHigh
	res.ResIdLow, res.ResIdHigh = binary.BigEndian.Uint64(id[0:8]), binary.BigEndian.Uint64(id[8:16])
	return nil
}

// get uuidv7 from reqbase
func GetUuidV7(req *oidcapiv1.Reqbase) uuid.UUID {
	id := uuid.UUID{}
	binary.BigEndian.PutUint64(id[0:8], req.ReqIdLow)
	binary.BigEndian.PutUint64(id[8:16], req.ReqIdHigh)
	return id
}
