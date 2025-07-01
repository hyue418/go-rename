package core

import (
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
)

// RenameImageAndVideo 根据拍摄时间重命名图片/视频文件
type RenameImageAndVideo struct {
	MatchFailureHandlerType int // 匹配失败的处理方式
}

func NewRenameImageAndVideo(matchFailureHandlerType int) *RenameImageAndVideo {
	return &RenameImageAndVideo{MatchFailureHandlerType: matchFailureHandlerType}
}

// CountFiles 统计需要重命名的文件数量
func (r *RenameImageAndVideo) CountFiles(dir string) (int64, error) {
	var fileCount int64 = 0
	if err := CheckMediainfoCommandExists(); err != nil {
		return 0, err
	}
	// 遍历目录及其子目录
	if err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Dir(path) == filepath.Join(dir, UnknownDateDir) {
			return nil
		}
		// 只处理图片和视频类型文件，过滤掉隐藏文件
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
func (r *RenameImageAndVideo) Rename(dir string, bar *mpb.Bar) error {
	// 遍历目录及其子目录
	if err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Dir(path) == filepath.Join(dir, UnknownDateDir) {
			return nil
		}
		// 过滤掉隐藏文件
		if file.IsDir() || IsHiddenFile(file.Name()) {
			return nil
		}
		if err = RenameSingleImageOrVideo(path, file, r.MatchFailureHandlerType); err != nil {
			return err
		}
		bar.Increment()
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// RenameSingleImageOrVideo 重命名单个图片/视频文件
func RenameSingleImageOrVideo(path string, file os.FileInfo, matchFailureHandlerType int) error {
	if IsImage(path) {
		return RenameSingleImage(path, file, matchFailureHandlerType)
	}
	if IsVideo(path) {
		return RenameSingleVideo(path, file, matchFailureHandlerType)
	}
	return nil
}
