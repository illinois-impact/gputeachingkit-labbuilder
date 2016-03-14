package main

import (
	"fmt"

	"github.com/facebookgo/stack"
	"github.com/fatih/color"
	"gitlab.com/abduld/wgx-pandoc/cmd"
)

func main() {

	defer func() {
		var (
			logFmt   = "\n[%s] %v \n\nStack Trace:\n============\n\n%s\n\n"
			titleClr = color.New(color.Bold, color.FgRed).SprintFunc()
		)
		if err := recover(); err != nil {
			frames := stack.Callers(4)
			fmt.Printf(logFmt, titleClr("PANIC"), err, frames.String())
		}
	}()

	cmd.Execute()
}
