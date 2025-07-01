package core

import (
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
)

// RenameFileByHash 根据文件hash重命名图片/视频文件
type RenameFileByHash struct {
}

func NewRenameFileByHash() *RenameFileByHash {
	return &RenameFileByHash{}
}

// CountFiles 统计需要重命名的文件数量
func (r *RenameFileByHash) CountFiles(dir string) (int64, error) {
	var fileCount int64 = 0
	// 遍历目录及其子目录
	if err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理图片和视频，过滤掉隐藏文件
		if (!IsImage(path) && !IsVideo(path)) || file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		fileCount++
		return nil
	}); err != nil {
		return 0, err
	}
	return fileCount, nil
}

// Rename 重命名
func (r *RenameFileByHash) Rename(dir string, bar *mpb.Bar) error {
	// 遍历目录及其子目录
	return filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理图片和视频，过滤掉隐藏文件
		if (!IsImage(path) && !IsVideo(path)) || file.IsDir() || IsHiddenFile(file.Name()) {
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
