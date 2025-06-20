package common

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	"github.com/fatih/color"
	"time"
)

const (
	DividingLine = "======================================================"
)

func Head() {
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

// PrintDividingLine 打印分割线
func PrintDividingLine() {
	fmt.Println("\n" + DividingLine + "\n")
}

// PrintError 打印错误
func PrintError(content string) {
	PrintDividingLine()
	color.New(color.FgRed).Add(color.Bold).Println(content + "\n")
}

// FormatDate 格式化日期
func FormatDate(time string) string {
	return carbon.Parse(time).Layout("20060102_150405")
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
