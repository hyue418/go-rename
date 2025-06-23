package core

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"hyue418/go-rename/common"
	"os"
)

// UnknownDateDir 未知拍摄日期文件夹
const UnknownDateDir = "unknown-date"

const (
	RenameTypeImage         = "rename-type-image"           // 重命名图片
	RenameTypeVideo         = "rename-type-video"           // 重命名视频
	RenameTypeImageAndVideo = "rename-type-image-and-video" // 重命名图片+视频
	RenameTypeFileByHash    = "rename-type-file-by-hash"    // 根据md5重命名(用于文件去重)
)

// RenameTypeMap 编号与重命名类型映射
var RenameTypeMap = map[int]string{
	1:  RenameTypeImage,
	2:  RenameTypeVideo,
	3:  RenameTypeImageAndVideo,
	99: RenameTypeFileByHash,
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

// RenameStrategy 重命名策略器
type RenameStrategy interface {
	CountFiles(dir string) (int64, error)
	Rename(dir string, bar *mpb.Bar) error
}

// Execute 执行
func Execute() {
	var dir, renameType string
	var matchFailureHandlerType int
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			var renameTypeNum, matchFailureHandlerTypeNum int
			fmt.Print("请输入要处理的目录路径：")
			_, err := fmt.Scanln(&dir)
			if dir == "" {
				common.PrintError("输入错误：请输入正确的目录路径")
				os.Exit(1)
			}
			if err != nil {
				common.PrintError("输入错误:" + err.Error())
				os.Exit(1)
			}
			if info, err := os.Stat(dir); os.IsNotExist(err) {
				common.PrintError("输入错误：目录不存在")
				os.Exit(1)
			} else if !info.IsDir() {
				common.PrintError("输入错误：请输入正确的目录路径")
				os.Exit(1)
			}
			common.PrintDividingLine()
			color.New(color.FgBlue).Add(color.Bold).Println("【操作类型】")
			fmt.Print(
				"1.根据拍摄时间重命名图片文件\n" +
					"2.根据拍摄时间重命名视频文件\n" +
					"3.根据拍摄时间重命名图片+视频文件\n" +
					"99.根据文件md5重命名(用于文件去重)\n" +
					"\n请输入编号:")
			_, err = fmt.Scanln(&renameTypeNum)
			if err != nil {
				common.PrintError("输入错误:" + err.Error())
				os.Exit(1)
			}
			renameType = RenameTypeMap[renameTypeNum]
			if renameType == "" {
				common.PrintError("输入错误：请输入正确的编号")
				os.Exit(1)
			}
			if renameType != RenameTypeFileByHash {
				common.PrintDividingLine()
				color.New(color.FgBlue).Add(color.Bold).Println("【部分文件可能没有拍摄日期,想如何处理?】")
				fmt.Print(
					"1.忽略此类文件,不处理\n" +
						"2.使用文件的创建时间/修改时间(哪个时间早用哪个)替代拍摄时间\n" +
						"3.统一将这部分文件移至unknown-date文件夹,不修改文件名\n" +
						"\n请输入编号:")
				_, err = fmt.Scanln(&matchFailureHandlerTypeNum)
				if err != nil {
					common.PrintError("输入错误:" + err.Error())
					os.Exit(1)
				}
				matchFailureHandlerType = MatchFailureHandlerTypeMap[matchFailureHandlerTypeNum]
				if matchFailureHandlerType == 0 {
					common.PrintError("输入错误：请输入正确的编号")
					os.Exit(1)
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
	p := mpb.New(mpb.WithWidth(64))
	bar := p.New(fileCount,
		mpb.BarStyle().Lbound("").Filler("█").Tip("█").Padding("░").Rbound(""),
		mpb.PrependDecorators(
			decor.Name("处理进度", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO), "done"),
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
	color.New(color.FgGreen).Add(color.Bold).Println("\n处理完成")
}
