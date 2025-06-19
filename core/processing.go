package core

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UnknownDateDir 未知拍摄日期文件夹
var UnknownDateDir = "unknown-date"

var (
	RenameTypeFileByHash              = "file-by-hash"                    // 根据md5重命名(用于文件去重)
	RenameTypePhotoByExifDate         = "photo-by-shooting-time"          // 根据照片的拍摄时间重命名，没有拍摄时间的文件移至unknown-date文件夹
	RenameTypePhotoByExifDatePriority = "photo-by-shooting-time-priority" // 优先根据照片的拍摄时间重命名，没有拍摄时间的按文件创建时间命名
	//RenameTypePhotoAndVideo           = "photo-and-video-by-shooting-time" // 根据照片和视频的拍摄时间重命名，没有拍摄时间的文件移至unknown-date文件夹
)

// GetRenameByDateFileCount 获取需要根据拍摄时间重命名的文件数量
func GetRenameByDateFileCount(dir string) (int64, error) {
	var fileCount int64 = 0
	// 遍历目录及其子目录
	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Dir(path) == filepath.Join(dir, UnknownDateDir) {
			return nil
		}
		// 只处理文件不处理目录，过滤掉隐藏文件
		if file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		fileCount++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return fileCount, nil
}

// RenameByDate 重命名文件，根据文件的拍摄时间生成新文件名
func RenameByDate(dir string, bar *mpb.Bar, isOnlyExifDate bool) error {
	// 遍历目录及其子目录
	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Dir(path) == filepath.Join(dir, UnknownDateDir) {
			return nil
		}
		// 只处理文件不处理目录，过滤掉隐藏文件
		if file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		originalTime, err := GetOriginalTime(path)
		if err != nil {
			return err
		}
		// 没有EXIF拍摄日期
		if originalTime == "" {
			if isOnlyExifDate {
				// 没有EXIF的文件移至unknown-date文件夹
				targetDir := filepath.Join(dir, UnknownDateDir)
				// 检查目标目录是否存在
				if _, err := os.Stat(targetDir); os.IsNotExist(err) {
					if err = os.Mkdir(targetDir, 0755); err != nil {
						return err
					}
				}
				// 构建目标文件的完整路径
				targetPath := filepath.Join(targetDir, filepath.Base(path))
				// 移动文件
				targetPath, err = renameWithConflictResolution(path, targetPath)
				if err != nil {
					fmt.Printf("Error move %s to %s: %v\n", path, targetPath, err)
				}
				bar.Increment()
				return nil
			}
			// 没有EXIF的按文件创建时间命名
			creationTime, _ := GetFileCreationTime(path)
			originalTime = creationTime
		}
		newFilePath := filepath.Join(filepath.Dir(path), GetDateFileName(originalTime, file.Name()))
		// 重命名文件
		newFilePath, err = renameWithConflictResolution(path, newFilePath)
		if err != nil {
			fmt.Printf("Error renaming %s to %s: %v\n", path, newFilePath, err)
		}
		bar.Increment()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetRenameByHashFileCount 获取需要根据文件MD5重命名的文件数量
func GetRenameByHashFileCount(dir string) (int64, error) {
	var fileCount int64 = 0
	// 遍历目录及其子目录
	if err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理文件不处理目录，过滤掉隐藏文件
		if file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		fileCount++
		return nil
	}); err != nil {
		return 0, err
	}
	return fileCount, nil
}

// RenameByHash 重命名文件，根据文件的 MD5 哈希值生成新文件名
func RenameByHash(dir string, bar *mpb.Bar) error {
	// 遍历目录及其子目录
	return filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理文件不处理目录，过滤掉隐藏文件
		if file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		newFileName, renameErr := GetNameByFileHash(path, file)
		if renameErr != nil {
			fmt.Printf("Error renaming %s: %v\n", path, renameErr)
			return nil
		}
		// 重命名文件
		err = os.Rename(path, newFileName)
		if err != nil {
			fmt.Printf("Error renaming %s to %s: %v\n", path, newFileName, err)
		}
		bar.Increment()
		return nil
	})
}

// renameWithConflictResolution 封装文件重命名，处理重名情况
func renameWithConflictResolution(oldPath, newPath string) (string, error) {
	// 检查目标文件是否是文件自身
	if oldPath == newPath {
		return oldPath, nil
	}
	// 检查目标文件是否存在，存在则加后缀
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		base := strings.TrimSuffix(newPath, GetExt(newPath))
		ext := GetExt(newPath)
		counter := 1
		for {
			newPath = fmt.Sprintf("%s_%d%s", base, counter, ext)
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				break
			}
			counter++
		}
	}
	// 执行重命名
	return newPath, os.Rename(oldPath, newPath)
}

// formatDate 格式化日期
func formatDate(time string) string {
	return carbon.Parse(time).Layout("20060102_150405")
}

// GetDateFileName 获取带日期的文件名(含后缀)
func GetDateFileName(date, fileName string) string {
	prefix := "FILE"
	if IsImage(fileName) {
		prefix = "IMG"
	} else if IsVideo(fileName) {
		prefix = "VID"
	}
	return fmt.Sprintf("%s_%s%s", prefix, formatDate(date), GetExt(fileName))
}

// minDateTime 返回时间字符串中的最小值，忽略空字符串
func minDateTime(times []string) (string, error) {
	var minTime *time.Time
	for _, dt := range times {
		if dt == "" {
			continue
		}
		parsedTime, err := time.Parse("2006-01-02 15:04:05", dt)
		if err != nil {
			return "", err
		}
		if minTime == nil || parsedTime.Before(*minTime) {
			minTime = &parsedTime
		}
	}
	if minTime != nil {
		return minTime.Format("2006-01-02 15:04:05"), nil
	}
	return "", nil
}
