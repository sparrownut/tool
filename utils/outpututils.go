package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"tool/Global"
)

func Printsuc(text string, args ...any) {
	c := color.New(color.FgHiGreen, color.Bold)
	_, _ = c.Printf(text+"\n", args...)
	if Global.OUTPUTFILE != "" {
		// 设置了就写入
		logFile, err := os.OpenFile(Global.OUTPUTFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			panic(err)
		}
		_, _ = logFile.WriteString(fmt.Sprintf(text))
	}

}
func Printerr(text string, args ...any) {
	c := color.New(color.FgHiRed)
	_, _ = c.Printf(text+"\n", args...)
}
func Printminfo(text string, args ...any) {
	c := color.New(color.FgYellow)
	_, _ = c.Printf(text+"\n", args...)
}
func Printhinfo(text string, args ...any) {
	c := color.New(color.FgHiYellow, color.Bold)
	_, _ = c.Printf(text+"\n", args...)
}
func Printcritical(text string, args ...any) {
	c := color.New(color.FgHiBlue, color.BgHiRed, color.Bold)
	_, _ = c.Printf(text+"\n", args...)
}
