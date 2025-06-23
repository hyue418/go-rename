package core

import (
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
)

// RenameImage 根据拍摄时间重命名图片文件
type RenameImage struct {
	MatchFailureHandlerType int // 匹配失败的处理方式
}

func NewRenameImage(matchFailureHandlerType int) *RenameImage {
	return &RenameImage{MatchFailureHandlerType: matchFailureHandlerType}
}

// CountFiles 统计需要重命名的文件数量
func (r *RenameImage) CountFiles(dir string) (int64, error) {
	var fileCount int64 = 0
	// 遍历目录及其子目录
	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Dir(path) == filepath.Join(dir, UnknownDateDir) {
			return nil
		}
		// 只处理图片类型文件，过滤掉隐藏文件
		if !IsImage(path) || file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		fmt.Println(path)
		fileCount++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return fileCount, nil
}

// Rename 重命名
func (r *RenameImage) Rename(dir string, bar *mpb.Bar) error {
	// 遍历目录及其子目录
	if err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Dir(path) == filepath.Join(dir, UnknownDateDir) {
			return nil
		}
		// 只处理图片类型文件，过滤掉隐藏文件
		if !IsImage(path) || file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		if err = RenameSingleImage(path, file, r.MatchFailureHandlerType); err != nil {
			return err
		}
		bar.Increment()
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// RenameSingleImage 重命名单张图片
func RenameSingleImage(path string, file os.FileInfo, matchFailureHandlerType int) error {
	if !IsImage(path) {
		return nil
	}
	originalTime, err := GetOriginalTime(path)
	if err != nil {
		return err
	}
	// 没有EXIF拍摄日期
	if originalTime == "" {
		switch matchFailureHandlerType {
		case MatchFailureHandlerTypeIgnore:
			return nil
		case MatchFailureHandlerTypeMoveToUnknownDateDir:
			// 没有EXIF的文件移至unknown-date文件夹
			targetDir := filepath.Join(filepath.Dir(path), UnknownDateDir)
			// 检查目标目录是否存在
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				if err = os.Mkdir(targetDir, 0755); err != nil {
					return err
				}
			}
			// 构建目标文件的完整路径
			targetPath := filepath.Join(targetDir, filepath.Base(path))
			// 移动文件
			targetPath, err = RenameWithConflictResolution(path, targetPath)
			if err != nil {
				fmt.Printf("Error move %s to %s: %v\n", path, targetPath, err)
			}
			return nil
		case MatchFailureHandlerTypeUseFileCreationTime:
			// 没有EXIF的按文件创建时间命名
			creationTime, _ := GetFileCreationTime(path)
			originalTime = creationTime
		}
	}
	newFilePath := filepath.Join(filepath.Dir(path), GetDateFileName(originalTime, file.Name()))
	// 重命名文件
	if newFilePath, err = RenameWithConflictResolution(path, newFilePath); err != nil {
		fmt.Printf("Error renaming %s to %s: %v\n", path, newFilePath, err)
	}
	return nil
}
