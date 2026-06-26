package kits

import (
	"runtime"
)

type _config struct {
}

var Config = _config{}

// GetMainPackageFile climbs the call stack to find the entry main function's directory.
func (c _config) GetMainPackageFile() string {
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

// GetCallSourceFile
func (c _config) GetCallSourceFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file
}
