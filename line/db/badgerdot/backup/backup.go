package backup

// 1. 保留两天的数据，每天一个全备份，每小时一个差异备份. 程序做好后，启动系统的定时间任务，每小时运行一次。 也有可能有手动的备份
// 2. 程序由操作系统的任务计划每小时启动一次，启动时间为整点过一分钟
// 3. 压缩备份文件
// 4. 备份文件名：backup_0_865.gzip, backup_865_996.gzip, 其中0_865为备份的开始与结束since, since是由badger.backup时返回的值，
//   从0开始的是为full backup file,下一个文件的开始是上一个文件的结束，这样的文件就是一组。如果备份中有不连续的文件，说明备份文件有缺失
//   如果一个备份文件不属于任何一个full backup, 那么当前备份目录也是不正常的
// 5. 在list backup file 时，full backup file 分组，组内是按照结束since从小到大排序，在组之间是按照full backup file的ts（file.ModTime）从小到大排序
import (
	"cmp"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/scryg/sutils/sfile"
)

type DbBackup struct {
	db          *badger.DB
	stopEvent   chan struct{}
	backupStart atomic.Bool
	backupPath  string
	Logger      *dot.LoggerType
}

type DbBackupFile struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// just file name of the backup file, no path
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// id that belong to which full backup file
	BelongTo   string `protobuf:"bytes,3,opt,name=belong_to,json=belongTo,proto3" json:"belong_to,omitempty"`
	Length     int64  `protobuf:"varint,4,opt,name=length,proto3" json:"length,omitempty"`
	Ts         int64  `protobuf:"varint,5,opt,name=ts,proto3" json:"ts,omitempty"`
	StartSince uint64 `protobuf:"varint,6,opt,name=start_since,json=startSince,proto3" json:"start_since,omitempty"`
	EndSince   uint64 `protobuf:"varint,7,opt,name=end_since,json=endSince,proto3" json:"end_since,omitempty"`
}
type DbBackupFolder struct {
	// fold name
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// the real file count of the backup folder, it is less equal len(files);
	FileCount int32 `protobuf:"varint,2,opt,name=file_count,json=fileCount,proto3" json:"file_count,omitempty"`
	// the back files are continuous
	Valid bool            `protobuf:"varint,3,opt,name=valid,proto3" json:"valid,omitempty"`
	Files []*DbBackupFile `protobuf:"bytes,4,rep,name=files,proto3" json:"files,omitempty"`
}

func (p *DbBackup) Backup(full bool) (*DbBackupFile, error) {
	now := time.Now()
	if full {
		return p.backupFull(now)
	} else {
		return p.backupDiff(now)
	}
}

// RestoreFile restores the database from a backup file.
// fileName is relative backupPath.
func (p *DbBackup) RestoreFile(fileId string) error {
	fullFile := filepath.Join(p.backupPath, fileId)
	root, err := os.OpenRoot(p.backupPath)
	if err != nil {
		p.Logger.Error().AnErr("cant open backup root", err).Send()
		return err
	}
	file, err := root.Open(fullFile)
	if err != nil {
		p.Logger.Info().AnErr("cant open the backup file", err).Send()
		return err
	}
	defer func() {
		if file != nil {
			err := file.Close()
			p.Logger.Error().AnErr("cant close the backup file", err).Send()
			file = nil
		}
	}()

	gzipFile, err := gzip.NewReader(file)
	if err != nil {
		p.Logger.Info().AnErr("cant create gzip reader", err).Send()
		return err
	}
	defer func() {
		if gzipFile != nil {
			err := gzipFile.Close()
			p.Logger.Error().AnErr("cant cloase the gzip file", err).Send()
			gzipFile = nil
		}
	}()

	err = p.db.Load(gzipFile, 16)
	if err != nil {
		p.Logger.Info().AnErr("cant load gzip file", err).Send()
		return err
	}

	return nil
}

func (p *DbBackup) RestoreFolder(folderId string) error {
	backupFolder, err := p.listBackupFilesInFolder(folderId)
	if err != nil {
		return err
	}
	if backupFolder == nil || len(backupFolder.Files) < 1 {
		return fmt.Errorf("the day folder is empty or cant read")
	}
	for i := range backupFolder.Files {
		err := p.RemoveFile(backupFolder.Files[i].Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *DbBackup) DropAllData() error {
	err := p.db.DropAll()
	if err != nil {
		p.Logger.Info().AnErr("cant drop all data", err).Send()
		return err
	}
	return nil
}

// RemoveFolder removes a backup folder by its ID
// the folderId must be in backupPath()
func (p *DbBackup) RemoveFolder(folderId string) error {
	if len(folderId) == 0 {
		return nil
	}
	absPath, err := filepath.Abs(filepath.Join(p.backupPath, folderId))
	if err != nil {
		return err
	}
	// security check: ensure the folder path is within the backup path
	if strings.HasPrefix(absPath, p.backupPath) {
		return os.RemoveAll(absPath)
	} else {
		p.Logger.Error().Msgf("folder path is outside the backup path: %s", absPath)
		return fmt.Errorf("folder path is outside the backup path")
	}
}

// RemoveFile removes a backup file by its ID
// the fileId must be in backupPath()
func (p *DbBackup) RemoveFile(fileId string) error {
	if len(fileId) == 0 {
		return nil
	}
	absFile, err := filepath.Abs(filepath.Join(p.backupPath, fileId))
	if err != nil {
		return err
	}
	// security check: ensure the file path is within the backup path
	if !strings.HasPrefix(absFile, p.backupPath) {
		p.Logger.Error().Msgf("file path is outside the backup path: %s", absFile)
		return fmt.Errorf("file path is outside the backup path")
	}
	return os.Remove(absFile)
}

func (p *DbBackup) List() ([]*DbBackupFolder, error) {
	backupFolder := []*DbBackupFolder{}
	entries, err := os.ReadDir(p.backupPath)
	if err != nil {
		p.Logger.Error().AnErr("cant read backup path", err).Send()
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			folder, err := p.listBackupFilesInFolder(entry.Name())
			if err != nil || folder == nil {
				continue
			}
			backupFolder = append(backupFolder, folder)
		}
	}
	return backupFolder, nil
}

func (p *DbBackup) listDayFolders() ([]*DbBackupFolder, error) {
	backupFolders := []*DbBackupFolder{}
	entries, err := os.ReadDir(p.backupPath)
	if err != nil {
		p.Logger.Error().AnErr("cant read backup path", err).Send()
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			backupFolder := DbBackupFolder{
				Id: entry.Name(),
			}
			backupFolders = append(backupFolders, &backupFolder)
		}
	}
	slices.SortFunc(backupFolders, func(a, b *DbBackupFolder) int {
		return cmp.Compare(a.Id, b.Id)
	})
	return backupFolders, nil
}

// the files in the folder are sorted
func (p *DbBackup) ValidFolder(folder *DbBackupFolder) bool {
	start := uint64(0)
	for _, file := range folder.Files {
		if file.StartSince == 0 {
			start = file.EndSince
			continue
		}
		if start != file.StartSince {
			folder.Valid = false
			return false
		}
		start = file.EndSince
	}
	folder.Valid = true
	return true
}

// group by full backup file, sort by ts(file.ModTime)
// in a group, sort by endSince
func (p *DbBackup) listBackupFilesInFolder(folderIdOr string) (*DbBackupFolder, error) {
	folderId, absFolder := p.parseBackupFolderId(folderIdOr)
	entries, err := os.ReadDir(absFolder)
	if err != nil {
		p.Logger.Error().AnErr("cant read folder path", err).Send()
		return nil, err
	}
	backupFolder := DbBackupFolder{
		Id:    folderId,
		Files: []*DbBackupFile{},
	}
	for _, subEntry := range entries {
		if !subEntry.IsDir() {
			info, err := subEntry.Info()
			if err != nil {
				p.Logger.Error().AnErr("cant read file info", err).Send()
				continue
			}
			name := subEntry.Name()
			f, err := parseBackupFileName(subEntry.Name())
			if err != nil {
				continue
			}
			file := DbBackupFile{
				Id:         filepath.Join(backupFolder.Id, name),
				Name:       name,
				Length:     info.Size(),
				Ts:         info.ModTime().Unix(),
				StartSince: f.start,
				EndSince:   f.end,
			}
			backupFolder.Files = append(backupFolder.Files, &file)
		}
	}
	backupFolder.FileCount = int32(len(backupFolder.Files)) // #nosec G115
	fullBackupFiles := make([]*DbBackupFile, 0, 8)
	{ // has more than one full backup file(the since start from zero)
		for _, file := range backupFolder.Files {
			if file.StartSince == 0 {
				fullBackupFiles = append(fullBackupFiles, file)
			}
		}
		slices.SortFunc(fullBackupFiles, func(a, b *DbBackupFile) int {
			return cmp.Compare(a.Ts, b.Ts)
		})
	}
	if len(fullBackupFiles) > 0 {
		groupFiles := make([][]*DbBackupFile, 0, len(fullBackupFiles))
		for _, file := range fullBackupFiles {
			file.BelongTo = file.Id
			lastGroup := []*DbBackupFile{file}
			start := file.EndSince
			for _, f := range backupFolder.Files {
				if f.StartSince != 0 && f.StartSince == start {
					f.BelongTo = file.Id
					lastGroup = append(lastGroup, f)
					start = f.EndSince
				}
			}
			slices.SortFunc(lastGroup, func(a, b *DbBackupFile) int {
				return cmp.Compare(a.EndSince, b.EndSince)
			})
			groupFiles = append(groupFiles, lastGroup)
		}
		tempFiles := make([]*DbBackupFile, 0, len(backupFolder.Files))
		for _, group := range groupFiles {
			tempFiles = append(tempFiles, group...)
		}
		// add the file not belong to full backup
		for _, file := range backupFolder.Files {
			if len(file.BelongTo) < 1 {
				tempFiles = append(tempFiles, file)
			}
		}
		backupFolder.Files = tempFiles

	}
	return &backupFolder, nil
}

func (p *DbBackup) start() {
	// check the backup path
	p.Logger.Info().Msgf("backup path: %s", p.backupPath)
	// 一小时运行一次，运行时间在整点过1分种，
	// 每次运行时，检查当天的全量是否备份，如果没有就进行全量备份
	// 如果全量已备份，就是进行差异备份
	go func() {
		for {
			now := time.Now()
			if now.Minute() > 0 && now.Minute() < 5 {
				if now.Hour() == 0 {
					_, err := p.backupFull(now)
					if err != nil {
						p.Logger.Error().AnErr("backup full error: ", err).Send()
					}
					err = p.removeOldBackup()
					if err != nil {
						p.Logger.Error().AnErr("remove old backup error: ", err).Send()
					}
				} else {
					_, err := p.backupDiff(now)
					if err != nil {
						p.Logger.Error().AnErr("backup diff error:", err).Send()
					}
				}
			}
			now = time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 1, 0, 0, now.Location())

			select {
			case <-p.stopEvent:
				return
			case <-time.After(time.Until(next)):
			}
		}
	}()
}
func (p *DbBackup) stop() {
	close(p.stopEvent)
}

func (p *DbBackup) backupDiff(now time.Time) (*DbBackupFile, error) {
	dayFolder := filepath.Join(p.backupPath, makeDayFolder(now))
	if !sfile.ExistDir(dayFolder) {
		if err := os.MkdirAll(dayFolder, 0750); err != nil {
			p.Logger.Error().AnErr("cant create day folder to backup", err).Send()
			return nil, err
		}
	}
	// get since
	var since uint64 = 0
	{
		folder, err := p.listBackupFilesInFolder(dayFolder)
		if err != nil {
			p.Logger.Error().AnErr("cant read day folder to backup", err).Send()
			return nil, err
		}
		if len(folder.Files) > 0 {
			latest := folder.Files[len(folder.Files)-1]
			since = latest.EndSince
		} // else no files found, go to full backup

	}
	if since == 0 {
		return p.backupFull(now)
	}
	return p.backupFile(dayFolder, since)
}

func (p *DbBackup) backupFull(now time.Time) (*DbBackupFile, error) {
	dayFolder := filepath.Join(p.backupPath, makeDayFolder(now))
	if !sfile.ExistDir(dayFolder) {
		if err := os.MkdirAll(dayFolder, 0750); err != nil {
			p.Logger.Error().AnErr("cant create day folder to backup", err).Send()
			return nil, err
		}
	}
	return p.backupFile(dayFolder, 0)
}
func (p *DbBackup) removeOldBackup() error {
	folders, err := p.listDayFolders()
	if err != nil {
		return err
	}
	if len(folders) > 2 {
		removes := folders[:len(folders)-2]
		for _, folder := range removes {
			err = p.RemoveFolder(folder.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// backupFile
func (p *DbBackup) backupFile(folder string, startSince uint64) (*DbBackupFile, error) {
	folderId, _ := p.parseBackupFolderId(folder)
	backupFile := makeTempBackupFileName(filepath.Join(p.backupPath, folderId))
	if sfile.ExistFile(backupFile) {
		err := os.Remove(backupFile)
		if err != nil {
			p.Logger.Error().AnErr("cant remove old backup", err).Send()
			return nil, err
		}
	}
	root, err := os.OpenRoot(p.backupPath)
	if err != nil {
		p.Logger.Error().AnErr("cant open backup root", err).Send()
		return nil, err
	}
	relRoot, err := filepath.Rel(p.backupPath, backupFile)
	if err != nil {
		p.Logger.Error().AnErr("cant rel backup root", err).Send()
		return nil, err
	}
	file, err := root.Create(relRoot)
	if err != nil {
		p.Logger.Error().AnErr("cant create backup file", err).Send()
		return nil, err
	}
	defer func() {
		if file != nil {
			err := file.Close()
			if err != nil {
				p.Logger.Error().AnErr("cant close backup file", err).Send()
			}
			file = nil
		}
	}()

	gzipFile := gzip.NewWriter(file)
	defer func() {
		if gzipFile != nil {
			err := gzipFile.Close()
			if err != nil {
				p.Logger.Error().AnErr("cant close backup gzip file", err).Send()
			}
			gzipFile = nil
		}
	}()

	endSince, err := p.db.Backup(gzipFile, startSince)
	if err != nil {
		p.Logger.Error().AnErr("cant backup db", err).Send()
		return nil, err
	}
	p.Logger.Info().Msgf("backup done, since_%d", endSince)
	// rename file with since
	{
		err := gzipFile.Close()
		gzipFile = nil
		if err != nil {
			p.Logger.Error().AnErr("cant close backup gzip file", err).Send()
			return nil, err
		}

		err = file.Close()
		file = nil
		if err != nil {
			p.Logger.Error().AnErr("cant close backup file", err).Send()
			return nil, err
		}
	}
	if endSince == 0 || endSince == startSince {
		p.Logger.Info().Msgf("backup not needed, since_%d, newSince_%d", startSince, endSince)
		err = os.Remove(backupFile)
		if err != nil {
			p.Logger.Error().AnErr("cant remove backup file", err).Send()
		}
		return nil, nil
	}
	name := makeBackupFileName(startSince, endSince)
	newFullFileName := filepath.Join(p.backupPath, folderId, name)
	if err := os.Rename(backupFile, newFullFileName); err != nil {
		p.Logger.Error().AnErr("cant rename backup file", err).Send()
		return nil, err
	}
	fileInfo, err := os.Stat(newFullFileName)
	if err != nil {
		p.Logger.Error().AnErr("cant get backup file info", err).Send()
		return nil, err
	}

	return &DbBackupFile{
		Id:         filepath.Join(folderId, name),
		Name:       name,
		Length:     fileInfo.Size(),
		Ts:         fileInfo.ModTime().Unix(),
		StartSince: startSince,
		EndSince:   endSince,
	}, nil
}

// 1 return folder id, 2 return full path
func (p *DbBackup) parseBackupFolderId(folder string) (string, string) {
	if filepath.IsAbs(folder) {
		return filepath.Base(folder), folder
	} else {
		return folder, filepath.Join(p.backupPath, folder)
	}
}

func NewDbBackup(config *badgerdot.BaderDbDotConfig, sconfig dot.SConfig, db *badgerdot.BadgerDbDot, logger *dot.LoggerType) (*DbBackup, func(), error) {
	backupPath, err := sconfig.FullPath(config.BackupPath)
	if err != nil {
		return nil, nil, err
	}
	back := &DbBackup{db: db.Db(), stopEvent: make(chan struct{}), backupPath: backupPath, backupStart: atomic.Bool{}, Logger: logger}
	back.backupStart.Store(false) // set it false
	back.start()
	return back, func() { back.stop() }, nil
}
