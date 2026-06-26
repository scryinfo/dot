package line

import (
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/scryinfo/dot/dot"
)

// VERSION_PKG :=github.com/scryinfo/dot/line
// COMMIT_MSG :=${shell git log -1 --pretty=%B}
// COMMIT_TIME :=${shell git show -s --format=%at}
// COMMIT_HASH :=${shell git rev-parse --short HEAD}
// BUILD_TIME :=${shell date +%s}

// build:
// 	${go} build -tags="release" -buildvcs=true -ldflags="-s -w \
// 	-X '${VERSION_PKG}.CommitMsg=${COMMIT_MSG}' \
// 	-X '${VERSION_PKG}.CommitTime=${COMMIT_TIME}' \
//  -X '${VERSION_PKG}.CommitHash=${COMMIT_HASH}' \
// 	-X '${VERSION_PKG}.BuildTime=${BUILD_TIME}'" -o dot_line main.go

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

func LogBuildInfo(info BuildInfo, logger *dot.LoggerType) {
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
	dot.Logger.Info().Msgf("commit msg: %s\ncommit time: %s\nbuild time: %s\ncommit hash: %s\nversion: %s\ngo version: %s\n", info.CommitMsg, info.CommitTime, info.BuildTime, info.CommitHash, info.Version, info.GoVersion)
}
