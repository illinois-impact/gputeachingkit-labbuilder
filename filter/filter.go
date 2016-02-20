package main

import (
	"encoding/json"
	"gitlab.com/abduld/wgx-pandoc/pkg"
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
	"os"

	"fmt"
	"github.com/facebookgo/stack"
	"github.com/fatih/color"
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

	var (
		format string
		doc    interface{}
	)

	decoder := json.NewDecoder(os.Stdin)
	decoder.Decode(&doc)
	if len(os.Args) > 1 {
		format = os.Args[1]
	} else {
		format = ""
	}

	meta := doc.([]interface{})[0].(map[string]interface{})["unMeta"]
	for _, filter := range pkg.Filters {
		doc = pf.Walk(doc, filter, format, meta)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(doc)
}
