package sconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
)

var OkGenerateConfig = fmt.Errorf("generate config ok")

// GenerateConfig generates a config file from the config struct
// config: the config struct
// rootConfig: the root config struct
// args: whether to parse command line arguments
// returns: whether to make config, error
func GenerateConfigWithArgs[T any](config dot.SConfig, rootConfig *T) (*T, error) {
	generateConfig := false
	{
		// parse command line arguments "GenerateConfig", dont use the flag package
		if len(os.Args) > 1 {
			name := "GenerateConfig"
			for i, it := range os.Args[1:] {
				it = strings.TrimSpace(it)
				suf := strings.TrimPrefix(it, "-"+name+"=")
				if suf == it {
					suf = strings.TrimPrefix(it, "--"+name+"=")
				}
				if suf == it {
					if it == "-"+name || it == "--"+name {
						if i < len(os.Args)-1 {
							itNext := strings.TrimSpace(os.Args[i+1])
							if !strings.HasPrefix(itNext, "-") {
								suf = itNext
							} else {
								suf = ""
							}
						} else {
							// the last argument is "GenerateConfig",
							suf = ""
						}
					}
				}

				if suf == "" {
					// the argument is "GenerateConfig",
					// no value, -GenerateConfig=true or --GenerateConfig=true or -GenerateConfig true or --GenerateConfig true)
					generateConfig = true
					break
				} else if suf != it {
					// the argument is "GenerateConfig",
					// has value, parse it as bool
					v, err := strconv.ParseBool(suf)
					if err == nil {
						generateConfig = v
					} else {
						fmt.Printf("cant parse GenerateConfig value: %s, set the GenerateConfig to false, err: %v\n", suf, err)
						generateConfig = false
					}
					break
				} else {
					// dont find the argument GenerateConfig, go next
				}
			}
		}
	}

	if generateConfig {
		ok, err := GenerateConfigGo(config, rootConfig)
		if err != nil {
			return nil, err
		}
		if ok {
			return rootConfig, OkGenerateConfig
		}
	}
	return rootConfig, nil
}

// GenerateConfigGo generates a config file from the config struct
// config: the config struct
// rootConfig: the root config struct
// returns: whether to make config, error
func GenerateConfigGo[T any](config dot.SConfig, rootConfig *T) (bool, error) {
	err := kits.Config.GenerateConfig(config, rootConfig)
	if err != nil {
		// dont use the logger here, the logger is not initialized yet
		fmt.Printf("make config err: %v\n", err)
		return false, err
	} else {
		fmt.Println("make config success")
		return true, nil
	}
}
