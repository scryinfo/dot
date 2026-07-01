package notstopupgrade

import (
	"net"
)

type UpgraderListenerConfig struct {
	PidFile string `json:"pid_file" toml:"pid_file" yaml:"pid_file" mapstructure:"pid_file"`
	// sample: "localhost:8080"
	Addr string `json:"addr" toml:"addr" yaml:"addr" mapstructure:"addr"`
}

type UpgraderListener struct {
	Listener     net.Listener
	WaitUpgrader func() error
}
