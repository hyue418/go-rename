package core

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/djherbis/times"
	"github.com/dsoprea/go-exif/v3"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"hyue418/go-rename/common"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// IsImage 判断文件是否为图片
func IsImage(path string) bool {
	extensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".svg", ".heif", ".heic", ".avif", ".ico",
		".cur", ".pcx", ".nef", ".cr2", ".jfif",
	}
	return funk.ContainsString(extensions, GetExt(path))
}

// IsVideo 判断文件是否为视频
func IsVideo(path string) bool {
	extensions := []string{
		".mp4", ".avi", ".mov", ".wmv", ".mkv", ".flv", ".webm", ".mpeg", ".mpg", ".3gp", ".3g2", ".m4v",
		".ogg", ".ogv", ".rm", ".rmvb", ".asf", ".divx", ".xvid", ".vob", ".m2v", ".m4p", ".mxf", ".mts",
		".m2ts", ".ts", ".tp", ".trp", ".f4v", ".f4p", ".f4a", ".f4b", ".ogm", ".dv", ".nsv", ".qt", ".rm",
		".ram", ".swf", ".slp", ".tp", ".trp", ".ts", ".vro", ".divx", ".xvid", ".img", ".vob", ".ifo",
		".dat", ".pva", ".rec", ".thp", ".tod", ".wtv", ".wtv", ".tp", ".trp", ".tivo", ".vdr", ".m2p",
		".m1a", ".m1s", ".m2a", ".m2s", ".m2t", ".m2ts", ".mts", ".mod", ".tod", ".wtv",
	}
	return funk.ContainsString(extensions, GetExt(path))
}

// GetExt 获取文件扩展名
func GetExt(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

// GetNameByFileHash 以文件md5作为文件名
func GetNameByFileHash(filePath string, info os.FileInfo) (string, error) {
	hash, err := GetFileHash(filePath)
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(filePath), hash+GetExt(info.Name())), nil
}

// GetOriginalTime 获取图片原始拍摄时间
func GetOriginalTime(filePath string) (string, error) {
	dateTimeOriginal, _, _, err := GetExifTime(filePath)
	if err != nil {
		return "", err
	}
	return dateTimeOriginal, nil
}

// GetFileCreationTime 获取文件创建时间
func GetFileCreationTime(filePath string) (string, error) {
	// 获取文件时间信息
	fileTimes, err := times.Stat(filePath)
	if err != nil {
		return "", err
	}
	// 创建时间
	birthTime := fileTimes.BirthTime()
	// 修改时间
	modTime := fileTimes.ModTime()
	// 对比返回最小时间
	if birthTime.Before(modTime) {
		return birthTime.Format("2006-01-02 15:04:05"), nil
	}
	return modTime.Format("2006-01-02 15:04:05"), nil
}

// GetExifTime 获取exif时间
func GetExifTime(filePath string) (dateTimeOriginal, dateTimeDigitized, dateTime string, err error) {
	opt := exif.ScanOptions{}
	dt, err := exif.SearchFileAndExtractExif(filePath)
	if errors.Is(err, exif.ErrNoExif) {
		return "", "", "", nil
	} else if err != nil {
		return
	}
	ets, _, err := exif.GetFlatExifData(dt, &opt)
	if err != nil {
		return
	}
	for _, et := range ets {
		switch et.TagName {
		case "DateTimeOriginal":
			dateTimeOriginal, err = FormatExifTime(fmt.Sprintf("%s", et.Value))
			if err != nil {
				return "", "", "", err
			}
		case "DateTimeDigitized":
			dateTimeDigitized, err = FormatExifTime(fmt.Sprintf("%s", et.Value))
			if err != nil {
				return "", "", "", err
			}
		case "DateTime":
			dateTime, err = FormatExifTime(fmt.Sprintf("%s", et.Value))
			if err != nil {
				return "", "", "", err
			}
		default:
			continue
		}
	}
	return dateTimeOriginal, dateTimeDigitized, dateTime, nil
}

// GetVideoDate 获取视频文件拍摄日期
func GetVideoDate(filename string) (string, error) {
	cmd := exec.Command("mediainfo", "--Output=JSON", filename)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	data := strings.TrimSpace(out.String())
	if data == "" {
		return "", nil
	}
	// 拍摄时间
	recordedDate := gjson.Get(data, `media.track.#(@type="General").Recorded_Date`).String()
	if recordedDate != "" {
		// 解析时间字符串
		res, err := time.Parse("2006-01-02T15:04:05-0700", recordedDate)
		if err != nil {
			return "", err
		}
		return res.Format("2006-01-02 15:04:05"), nil
	}
	// 编码时间
	encodedDate := gjson.Get(data, `media.track.#(@type="General").Encoded_Date`).String()
	if encodedDate != "" {
		// 解析时间字符串
		res, err := time.Parse("2006-01-02 15:04:05 MST", encodedDate)
		if err != nil {
			return "", err
		}
		return res.Format("2006-01-02 15:04:05"), nil
	}
	return "", nil
}

// FormatExifTime 格式化Exif时间
func FormatExifTime(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	res, err := time.Parse("2006:01:02 15:04:05", input)
	if err != nil {
		return "", err
	}
	return res.Format("2006-01-02 15:04:05"), nil
}

// IsHiddenFile 是否为隐藏文件
func IsHiddenFile(fileName string) bool {
	return len(fileName) > 0 && fileName[0] == '.'
}

// GetFileTime 获取文件时间
func GetFileTime(filePath string) (dateTimeOriginal, dateTimeDigitized, dateTime, modificationTime string, err error) {
	dateTimeOriginal, dateTimeDigitized, dateTime, err = GetExifTime(filePath)
	if err != nil {
		return
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	return dateTimeOriginal, dateTimeDigitized, dateTime, fileInfo.ModTime().Format("2006-01-02 15:04:05"), nil
}

// GetFileHash 计算文件的 MD5 哈希值
func GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// GetDateFileName 获取带日期的文件名(含后缀名)
func GetDateFileName(date, fileName string) string {
	prefix := "FIL"
	if IsImage(fileName) {
		prefix = "IMG"
	} else if IsVideo(fileName) {
		prefix = "VID"
	}
	return fmt.Sprintf("%s_%s%s", prefix, common.FormatDate(date), GetExt(fileName))
}

// GetFileInfo 获取文件信息
func GetFileInfo(filePath string) error {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	// 打印文件信息
	fmt.Println("文件名:", fileInfo.Name())
	fmt.Println("文件大小:", fileInfo.Size(), "字节")
	fmt.Println("文件权限:", fileInfo.Mode())
	fmt.Println("最后修改时间:", fileInfo.ModTime().Format(time.RFC1123))
	fmt.Println("是否为目录:", fileInfo.IsDir())
	return nil
}

// RenameWithConflictResolution 封装文件重命名，处理重名情况
func RenameWithConflictResolution(oldPath, newPath string) (string, error) {
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
