//go:build windows
// +build windows

package ende

import (
	"os"
	"syscall"
)

func SendSignal(s syscall.Signal) error {
	p, err := os.FindProcess(os.Getpid())

	if err != nil {
		return err
	}

	return p.Signal(s)
}
