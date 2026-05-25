//go:build windows

package rocksdbdot

/*
#cgo CFLAGS: -ID:/lang/vcpkg/installed/x64-mingw-static/include
#cgo LDFLAGS: -LD:/lang/vcpkg/installed/x64-mingw-static/lib
#cgo LDFLAGS: -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd
#cgo LDFLAGS: -static -static-libgcc -static-libstdc++
*/
import "C"

// #cgo LDFLAGS: -L${ROCKSDB_LIB}
// #cgo CFLAGS: -I${ROCKSDB_INCLUDE}
