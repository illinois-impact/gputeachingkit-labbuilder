package main

import (
	"encoding/json"
	"gitlab.com/abduld/wgx-labpdf/pkg"
	pf "gitlab.com/abduld/wgx-labpdf/pkg/pandocfilter"
	"os"
)

func main() {
	decoder := json.NewDecoder(os.Stdin)
	var data []interface{}
	var doc interface{}
	decoder.Decode(&data)
	var format string
	if len(os.Args) > 1 {
		format = os.Args[1]
	} else {
		format = ""
	}
	doc = data
	for _, filter := range pkg.Filters {
		doc = pf.Walk(doc, filter, format, data[0].(map[string]interface{})["unMeta"])
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(doc)
}
