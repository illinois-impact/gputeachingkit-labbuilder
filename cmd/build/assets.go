//go:generate rice embed-go

package build

import "github.com/GeertJohan/go.rice"

var (
	box              = rice.MustFindBox("./../../_fixtures")
	texTemplate      = box.MustString("md.template")
	markdownTemplate = box.MustString("tex.template")

	latexTemplateResources map[string]resource
)

func init() {

	latexTemplateBox := rice.MustFindBox("./../../_fixtures/latex_template")
	getTexResource := func(filename string) resource {
		return resource{
			fileName: filename,
			content:  latexTemplateBox.MustString(filename),
		}
	}
	latexTemplateResources = map[string]resource{
		"structure.tex": getTexResource("structure.tex"),
		"ccicons.sty":   getTexResource("ccicons.sty"),
		"template.tex":  getTexResource("template.tex"),
	}
}
