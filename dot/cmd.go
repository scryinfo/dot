// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"flag"
	"log"
	"os"
)

// Cmd type general command line parameters
type Cmd struct {
	//ConfigPath config file path
	ConfigPath string
	//ConfigFile config file name without path
	ConfigFile string
}

// GCmd Global variables general command line parameters
var GCmd Cmd

// CmdParameterName Command line parameter name
type CmdParameterName string

func (c CmdParameterName) String() string {
	return string(c)
}

// Command line parameters
const (
	//Config file path, add Cmd to show relations with struct cmd
	CmdConfigPath CmdParameterName = "configpath"
	//Config file name without path
	CmdConfigFile CmdParameterName = "configfile"
)

// // FlagDefines General command parameter initialization
// func FlagDefines() {
// 	flag.StringVar(&GCmd.ConfigPath, CmdConfigPath.String(), "", "config path")
// 	flag.StringVar(&GCmd.ConfigFile, CmdConfigFile.String(), "", "config file, not include path")
// }

func init() {
	fs := flag.NewFlagSet("dot", flag.ContinueOnError)
	fs.StringVar(&GCmd.ConfigPath, CmdConfigPath.String(), "", "config path")
	fs.StringVar(&GCmd.ConfigFile, CmdConfigFile.String(), "", "config file, not include path")
	err := fs.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
