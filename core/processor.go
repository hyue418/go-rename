package core

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"hyue418/go-rename/common"
	"os"
	"sort"
)

// UnknownDateDir 未知拍摄日期文件夹
const UnknownDateDir = "unknown-date"

const (
	RenameTypeImage         = "rename-type-image"           // 重命名图片
	RenameTypeVideo         = "rename-type-video"           // 重命名视频
	RenameTypeImageAndVideo = "rename-type-image-and-video" // 重命名图片+视频
	RenameTypeFileByHash    = "rename-type-file-by-hash"    // 根据md5重命名(用于文件去重)
)

// RenameTypeNumberMap 编号与重命名类型映射
var RenameTypeNumberMap = map[int]string{
	1:  RenameTypeImage,
	2:  RenameTypeVideo,
	3:  RenameTypeImageAndVideo,
	99: RenameTypeFileByHash,
}

// RenameTypeTextMap 重命名类型的文本映射
var RenameTypeTextMap = map[string]string{
	RenameTypeImage:         "根据拍摄时间重命名图片文件",
	RenameTypeVideo:         "根据拍摄时间重命名视频文件",
	RenameTypeImageAndVideo: "根据拍摄时间重命名图片+视频文件",
	RenameTypeFileByHash:    "根据文件md5重命名(用于文件去重)",
}

// 日期获取失败的处理方式
const (
	MatchFailureHandlerTypeIgnore               = 1 + iota // 忽略
	MatchFailureHandlerTypeUseFileCreationTime             // 使用文件创建时间
	MatchFailureHandlerTypeMoveToUnknownDateDir            // 移至unknown-date文件夹
)

// MatchFailureHandlerTypeMap 编号与日期获取失败的处理方式映射
var MatchFailureHandlerTypeMap = map[int]int{
	1: MatchFailureHandlerTypeIgnore,
	2: MatchFailureHandlerTypeUseFileCreationTime,
	3: MatchFailureHandlerTypeMoveToUnknownDateDir,
}

var MatchFailureHandlerTypeTextMap = map[int]string{
	MatchFailureHandlerTypeIgnore:               "忽略此类文件,不处理",
	MatchFailureHandlerTypeUseFileCreationTime:  "使用文件的创建时间/修改时间(哪个时间早用哪个)替代拍摄时间",
	MatchFailureHandlerTypeMoveToUnknownDateDir: "统一将这部分文件移至unknown-date文件夹,不修改文件名",
}

// RenameStrategy 重命名策略器
type RenameStrategy interface {
	CountFiles(dir string) (int64, error)
	Rename(dir string, bar *mpb.Bar) error
}

// Execute 执行
func Execute() {
	var dir, renameType string
	var matchFailureHandlerType int
	var numbers []int
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			var renameTypeNum, matchFailureHandlerTypeNum int
			var inputPassed bool
			var err error
			for !inputPassed {
				fmt.Print("请输入要处理的目录路径：")
				_, err = fmt.Scanln(&dir)
				if err != nil || dir == "" {
					common.PrintError("输入错误：请输入正确的目录路径")
					continue
				}
				if info, err := os.Stat(dir); os.IsNotExist(err) {
					common.PrintError("输入错误：目录不存在")
					continue
				} else if !info.IsDir() {
					common.PrintError("输入错误：请输入正确的目录路径")
					continue
				}
				inputPassed = true
			}
			common.PrintDividingLine()
			color.New(color.FgBlue).Add(color.Bold).Println("【操作类型】")
			numbers = funk.Keys(RenameTypeNumberMap).([]int)
			sort.Ints(numbers)
			for _, v := range numbers {
				fmt.Printf("%d.%s\n", v, RenameTypeTextMap[RenameTypeNumberMap[v]])
			}
			fmt.Println()
			inputPassed = false
			for !inputPassed {
				fmt.Print("请输入编号:")
				_, err = fmt.Scanln(&renameTypeNum)
				if err != nil {
					common.PrintError("输入错误：请输入正确的编号")
					continue
				}
				renameType = RenameTypeNumberMap[renameTypeNum]
				if renameType == "" {
					common.PrintError("输入错误：请输入正确的编号")
					continue
				}
				inputPassed = true
			}
			common.PrintDividingLine()
			if renameType != RenameTypeFileByHash {
				color.New(color.FgBlue).Add(color.Bold).Println("【部分文件可能没有拍摄日期,想如何处理?】")
				numbers = funk.Keys(MatchFailureHandlerTypeMap).([]int)
				sort.Ints(numbers)
				for _, v := range numbers {
					fmt.Printf("%d.%s\n", v, MatchFailureHandlerTypeTextMap[v])
				}
				fmt.Println()
				inputPassed = false
				for !inputPassed {
					fmt.Print("请输入编号:")
					_, err = fmt.Scanln(&matchFailureHandlerTypeNum)
					if err != nil {
						common.PrintError("输入错误：请输入正确的编号")
						continue
					}
					matchFailureHandlerType = MatchFailureHandlerTypeMap[matchFailureHandlerTypeNum]
					if matchFailureHandlerType == 0 {
						common.PrintError("输入错误：请输入正确的编号")
						continue
					}
					inputPassed = true
				}
			}
			common.PrintDividingLine()
			color.New(color.FgBlue).Add(color.Bold).Println("【操作确认】")
			color.New().Add(color.Bold).Printf("处理的目录路径：")
			color.New().Add(color.Bold).Add(color.Underline).Printf("%s\n", dir)
			color.New().Add(color.Bold).Printf("处理方式：")
			color.New().Add(color.Bold).Add(color.Underline).Printf("%s\n", RenameTypeTextMap[renameType])
			switch renameType {
			case RenameTypeImage:
				color.New().Add(color.Bold).Printf("该目录下(包含子目录)所有符合条件的图片文件将重命名为[IMG_20250606_121601.xxx]的格式")
			case RenameTypeVideo:
				color.New().Add(color.Bold).Printf("该目录下(包含子目录)所有符合条件的视频文件将重命名为[VID_20250606_121601.xxx]的格式")
			case RenameTypeImageAndVideo:
				color.New().Add(color.Bold).Printf("该目录下(包含子目录)所有符合条件的图片文件和视频文件将重命名为[IMG_20250606_121601.xxx]和[VID_20250606_121601.xxx]的格式")
			}
			fmt.Print("\n")
			confirmType := 2
			var confirmText string
			for confirmType == 2 {
				fmt.Print("确认处理吗？请输入Y/N\n")
				_, err = fmt.Scanln(&confirmText)
				if err != nil {
					common.PrintError("输入错误，请输入Y/N")
					continue
				}
				if confirmText == "Y" || confirmText == "y" {
					confirmType = 1
				} else if confirmText == "N" || confirmText == "n" {
					common.PrintError("结束运行")
					os.Exit(1)
				} else {
					common.PrintError("输入错误：请输入Y/N")
					continue
				}
			}
		},
	}
	if err := cmd.Execute(); err != nil {
		common.PrintError(err.Error())
		os.Exit(1)
	}
	common.PrintDividingLine()
	color.New(color.FgBlue).Add(color.Bold).Println("正在统计文件数量,请稍后...")
	var renameStrategy RenameStrategy
	switch renameType {
	case RenameTypeFileByHash:
		renameStrategy = NewRenameFileByHash()
	case RenameTypeImage:
		renameStrategy = NewRenameImage(matchFailureHandlerType)
	case RenameTypeVideo:
		renameStrategy = NewRenameVideo(matchFailureHandlerType)
	case RenameTypeImageAndVideo:
		renameStrategy = NewRenameImageAndVideo(matchFailureHandlerType)
	default:
		return
	}
	fileCount, err := renameStrategy.CountFiles(dir)
	if err != nil {
		common.PrintError(err.Error())
		os.Exit(1)
	}
	color.New(color.FgBlue).Add(color.Bold).Println(fmt.Sprintf("共计%d个需处理的文件,开始进行处理\n", fileCount))
	color.New().Add(color.Bold).Println("处理进度")
	p := mpb.New(mpb.WithWidth(64))
	bar := p.New(fileCount,
		mpb.BarStyle().Lbound("").Filler("█").Tip("█").Padding("░").Rbound(""),
		mpb.PrependDecorators(
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO), "已完成"),
			decor.CountersNoUnit(" %d/%d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)
	go func() {
		if err = renameStrategy.Rename(dir, bar); err != nil {
			common.PrintError(err.Error())
			os.Exit(1)
		}
	}()
	p.Wait()
	color.New(color.FgGreen).Add(color.Bold).Println("\n=======================处理完成=======================")
}
