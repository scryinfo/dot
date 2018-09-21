package dot

import (
	"flag"
)

//Cmd 类型 通用的命令行参数
type Cmd struct {
	//ConfigPath 配置文件路径
	ConfigPath string
	//ConfigFile 配置文件名，不包含路径
	ConfigFile string
}

//GCmd 全局变量 通用的命令行参数
var GCmd Cmd

//CmdParameterName 命令行参数名字
type CmdParameterName string

func (c CmdParameterName) String() string {
	return string(c)
}	

//命令行参数
const (
	//配置文件路径，增加Cmd以示与struct cmd 相关
	CmdConfigPath CmdParameterName = "configpath"
	//配置文件名，不包含路径
	CmdConfigFile CmdParameterName = "configfile"
)

//CmdDefines 通用命令参数的初始化
func FlagDefines() {
	flag.StringVar(&GCmd.ConfigPath, CmdConfigPath.String(), "", "config path")
	flag.StringVar(&GCmd.ConfigPath, CmdConfigFile.String(), "", "config file, not include path")
}
