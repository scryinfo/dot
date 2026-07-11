package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

const (
	generateFileBaseNameExplain    = "file base name of 'private key' and 'ciphertext config'"
	plaintextConfigFileNameExplain = "plaintext config file name, default 'xxx_local.toml'"
)

var (
	generateFileBaseName    string
	plaintextConfigFileName string

	help                bool
	configFileExtension string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&help, "help", false, "this help")

	flag.StringVar(&generateFileBaseName, "n", "", generateFileBaseNameExplain)
	flag.StringVar(&generateFileBaseName, "name", "", generateFileBaseNameExplain)

	flag.StringVar(&plaintextConfigFileName, "c", "", plaintextConfigFileNameExplain)
	flag.StringVar(&plaintextConfigFileName, "config", "", plaintextConfigFileNameExplain)

	flag.Parse()

	if help {
		//flag.PrintDefaults 会分开打印'h'和'help'，所以这里自己写
		log.Println("\n> Options:\n\n" +
			"  -h --help:\n\tthis help\n" +
			"  -n --name:\n\t" + generateFileBaseNameExplain + "\n" +
			"  -c --config:\n\t" + plaintextConfigFileNameExplain)
		os.Exit(0)
	}

	/* set default value */

	if generateFileBaseName == "" { // 默认使用当前路径名作为生成文件名
		path, err := os.Getwd()
		if err != nil {
			log.Fatalln("get executable path failed, error: ", err)
		}

		index := strings.LastIndex(path, "/")
		generateFileBaseName = path[index+1:]
	}

	if plaintextConfigFileName == "" {
		plaintextConfigFileName = generateFileBaseName + "_local.toml" // 默认是toml格式，其他格式请显式提供
	}

	index := strings.LastIndex(plaintextConfigFileName, ".")
	configFileExtension = plaintextConfigFileName[index+1:]
}
