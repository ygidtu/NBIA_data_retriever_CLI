package main

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

// bytesBar 自定义的progress风格，用以展示文件的读写进度
func bytesBar(size int64, name string) *progressbar.ProgressBar {
	if len(name) > 50 {
		name = fmt.Sprintf("%s...", name[0:51])
	}

	return progressbar.NewOptions(int(size),
		progressbar.OptionSetWriter(ansi.NewAnsiStderr()), // 使用ansi防止奇怪的换行等显示错误
		progressbar.OptionUseANSICodes(true),              // avoid progressbar downsize error
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription(name),
		progressbar.OptionFullWidth(),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetWidth(10),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
