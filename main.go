package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"hyue418/go-rename/core"
	"os"
)

const (
	DividingLine = "======================================================"
)

var (
	dir        string // 目录路径
	renameType string // 重命名类型
)

func head() {
	fmt.Println("\n" +
		"   ____ _____     ________  ____  ____ _____ ___  ___ \n" +
		"  / __ `/ __ \\   / ___/ _ \\/ __ \\/ __ `/ __ `__ \\/ _ \\\n" +
		" / /_/ / /_/ /  / /  /  __/ / / / /_/ / / / / / /  __/\n" +
		" \\__, /\\____/  /_/   \\___/_/ /_/\\__,_/_/ /_/ /_/\\___/ \n" +
		"/____/\n" +
		"┌────────────────────────────────────────────────────┐\n" +
		"│        https://github.com/hyue418/go-rename        │\n" +
		"└────────────────────────────────────────────────────┘\n")
}

func main() {
	head()
	renameTypeMap := map[int]string{
		1: core.RenameTypeFileByHash,
		2: core.RenameTypePhotoByExifDate,
	}

	var cmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			var renameTypeNum int
			fmt.Print("请输入要处理的目录路径：")
			_, err := fmt.Scanln(&dir)
			if dir == "" {
				printError("输入错误：请输入正确的目录路径")
				os.Exit(1)
			}
			if err != nil {
				printError("输入错误:" + err.Error())
				os.Exit(1)
			}
			if info, err := os.Stat(dir); os.IsNotExist(err) {
				printError("输入错误：目录不存在")
				os.Exit(1)
			} else if !info.IsDir() {
				printError("输入错误：请输入正确的目录路径")
				os.Exit(1)
			}
			printDividingLine()
			color.New(color.FgBlue).Add(color.Bold).Println("重命名的类型")
			fmt.Print("1.根据md5重命名(用于文件去重)\n" +
				"2.根据文件EXIF拍摄日期重命名，没有EXIF的文件移至no-exif文件夹\n" +
				"3.优先根据EXIF拍摄日期重命名，没有EXIF的按文件创建时间命名\n" +
				"\n请输入编号:")
			_, err = fmt.Scanln(&renameTypeNum)
			if err != nil {
				printError("输入错误:" + err.Error())
				os.Exit(1)
			}
			renameType = renameTypeMap[renameTypeNum]
			if renameType == "" {
				printError("输入错误：请输入正确的编号")
				os.Exit(1)
			}
		},
	}
	if err := cmd.Execute(); err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	printDividingLine()
	// 获取数量
	var fileCount int64 = 0
	color.New(color.FgBlue).Add(color.Bold).Println("正在统计文件数量,请稍后...")
	switch renameType {
	case core.RenameTypeFileByHash:
		var err error
		fileCount, err = core.GetRenameByHashFileCount(dir)
		if err != nil {
			fmt.Println(err)
		}
	case core.RenameTypePhotoByExifDate:
		fallthrough
	case core.RenameTypePhotoByExifDatePriority:
		var err error
		fileCount, err = core.GetRenameByDateFileCount(dir)
		if err != nil {
			fmt.Println(err)
		}
	default:
		return
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
		switch renameType {
		case core.RenameTypeFileByHash:
			if err := core.RenameByHash(dir, bar); err != nil {
				fmt.Println(err)
			}
		case core.RenameTypePhotoByExifDate:
			if err := core.RenameByDate(dir, bar, true); err != nil {
				fmt.Println(err)
			}
		case core.RenameTypePhotoByExifDatePriority:
			if err := core.RenameByDate(dir, bar, false); err != nil {
				fmt.Println(err)
			}
		default:
			return
		}
	}()
	p.Wait()
	color.New(color.FgGreen).Add(color.Bold).Println("\n处理完成")
	return
}

// printDividingLine 打印分割线
func printDividingLine() {
	fmt.Println("\n" + DividingLine + "\n")
}

// printError 打印错误
func printError(content string) {
	printDividingLine()
	color.New(color.FgRed).Add(color.Bold).Println(content + "\n")
}
