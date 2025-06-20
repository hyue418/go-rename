package core

import (
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
)

// RenamePhoto 根据拍摄时间重命名照片
type RenamePhoto struct {
	MatchFailureHandlerType int // 匹配失败的处理方式
}

func NewRenamePhoto(matchFailureHandlerType int) *RenamePhoto {
	return &RenamePhoto{MatchFailureHandlerType: matchFailureHandlerType}
}

// CountFiles 统计需要重命名的文件数量
func (r *RenamePhoto) CountFiles(dir string) (int64, error) {
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
		fileCount++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return fileCount, nil
}

// Rename 重命名
func (r *RenamePhoto) Rename(dir string, bar *mpb.Bar) error {
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
		originalTime, err := GetOriginalTime(path)
		if err != nil {
			return err
		}
		// 没有EXIF拍摄日期
		if originalTime == "" {
			switch r.MatchFailureHandlerType {
			case MatchFailureHandlerTypeIgnore:
				// 忽略
				bar.Increment()
				return nil
			case MatchFailureHandlerTypeMoveToUnknownDateDir:
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
				targetPath, err = RenameWithConflictResolution(path, targetPath)
				if err != nil {
					fmt.Printf("Error move %s to %s: %v\n", path, targetPath, err)
				}
				bar.Increment()
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
		bar.Increment()
		return nil
	}); err != nil {
		return err
	}
	return nil
}
