package dot

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var (
	CommitMsg  = "dont get log message"
	CommitTime = "" // the parameter "-X"  only string
	BuildTime  = ""
	CommitHash = ""
	Version    = ""
	GoVersion  = ""
)

const timeFormat = "2006-01-02 15:04:05"

func GetBuildInfo() BuildInfo {
	info := BuildInfo{
		CommitMsg:  CommitMsg,
		CommitTime: CommitTime,
		BuildTime:  BuildTime,
		CommitHash: CommitHash,
		Version:    Version,
		GoVersion:  GoVersion,
	}
	if info.CommitTime != "" {
		ts, err := time.Parse(time.RFC3339, info.CommitTime)
		if err == nil {
			info.CommitTime = ts.Format(timeFormat)
		}
	}
	if info.BuildTime != "" {
		ts, err := time.Parse(time.RFC3339, info.BuildTime)
		if err == nil {
			info.BuildTime = ts.Format(timeFormat)
		}
	} else {
		exeFile, err := os.Executable()
		if err == nil {
			fileInfo, err := os.Stat(exeFile)
			if err == nil {
				info.BuildTime = fileInfo.ModTime().Format(timeFormat)
			}
		}
	}
	{
		readInfo, ok := debug.ReadBuildInfo()
		if ok {
			info.GoVersion = readInfo.GoVersion
			if info.Version == "" {
				if len(readInfo.Main.Version) > 0 {
					if !strings.Contains(readInfo.Main.Version, "devel") {
						info.Version = readInfo.Main.Version
					}
				}
			}
			if info.CommitHash == "" || info.CommitTime == "" {
				for _, setting := range readInfo.Settings {
					switch setting.Key {
					case "vcs.revision":
						if info.CommitHash == "" {
							info.CommitHash = setting.Value
						}
					case "vcs.time":
						if info.CommitTime == "" {
							ts, err := time.Parse(time.RFC3339, setting.Value)
							if err == nil {
								info.CommitTime = ts.Format(timeFormat)
							}
						}
					}
				}
			}
		}
	}
	return info
}

type BuildInfo struct {
	CommitMsg  string
	CommitTime string
	BuildTime  string
	CommitHash string
	Version    string
	GoVersion  string
}

func LogBuildInfo(info BuildInfo, logger *LoggerType) {
	if logger == nil {
		FmtBuildInfo(info)
	} else {
		logger.Info().Fields(map[string]interface{}{
			"commitMsg":  info.CommitMsg,
			"commitTime": info.CommitTime,
			"buildTime":  info.BuildTime,
			"commitHash": info.CommitHash,
			"version":    info.Version,
			"goVersion":  info.GoVersion,
		}).Msg("")
	}

}

func FmtBuildInfo(info BuildInfo) {
	fmt.Printf("commit msg: %s\ncommit time: %s\nbuild time: %s\ncommit hash: %s\nversion: %s\ngo version: %s\n", info.CommitMsg, info.CommitTime, info.BuildTime, info.CommitHash, info.Version, info.GoVersion)
}
