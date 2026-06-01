package backup

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	backupPrefix    = "backup_"
	backupExt       = ".gzip"
	dayFolderFormat = "2006-01-02"
)

type backupFile struct {
	name  string // no path, just the file name
	start uint64
	end   uint64
}

func makeTempBackupFileName(folder string) string {
	return filepath.Join(folder, fmt.Sprintf("___%d%s", time.Now().UnixNano(), backupExt))
}

func makeBackupFileName(start uint64, end uint64) string {
	return fmt.Sprintf("%s%d_%d%s", backupPrefix, start, end, backupExt)
}

func makeDayFolder(dateTime time.Time) string {
	return dateTime.Format(dayFolderFormat)
}

func parseBackupFileName(backupFileName string) (backupFile, error) {
	file := backupFile{}
	file.name = filepath.Base(backupFileName)
	st := strings.TrimPrefix(file.name, backupPrefix)
	st = strings.TrimSuffix(st, backupExt)
	sts := strings.Split(st, "_")
	if len(sts) != 2 {
		return file, fmt.Errorf("invalid backup file name: %s", backupFileName)
	}
	var err error
	file.start, err = strconv.ParseUint(sts[0], 10, 64)
	if err != nil {
		return file, fmt.Errorf("invalid backup file name: %s", backupFileName)
	}
	file.end, err = strconv.ParseUint(sts[1], 10, 64)
	if err != nil {
		return file, fmt.Errorf("invalid backup file name: %s", backupFileName)
	}
	return file, nil
}
