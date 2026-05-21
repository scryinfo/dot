package kits

import (
	"runtime"
)

type _config struct {
}

var Config = _config{}

// GetMainPackageDir climbs the call stack to find the entry main function's directory.
func (c _config) GetMainPackageDir() string {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(1, pcs)
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if frame.Function == "main.main" {
			return frame.File
		}
		if !more {
			break
		}
	}
	return ""
}
