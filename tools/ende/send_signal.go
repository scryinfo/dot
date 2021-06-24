// +build !windows

package ende

import (
	"os"
	"syscall"
)

func SendSignal(s syscall.Signal) error {
	return syscall.Kill(os.Getegid(), s)
}
