package dot

import (
	"flag"
)

type Cmd_ struct {
	Confpath string
	Conffile string
}

var Cmd Cmd_

func CmdDefines() {
	flag.StringVar(&Cmd.Confpath, "confpath", "", "config path")
	flag.StringVar(&Cmd.Confpath, "conffile", "", "config file, not include path")
}
