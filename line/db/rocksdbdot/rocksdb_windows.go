//go:build windows

package rocksdbdot

/*
// #cgo LDFLAGS: -L${ROCKSDB_LIB}
// #cgo CFLAGS: -I${ROCKSDB_INCLUDE}
#cgo LDFLAGS: -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd
#cgo LDFLAGS: -static -static-libgcc -static-libstdc++
*/
import "C"
