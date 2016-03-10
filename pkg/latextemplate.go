//go:generate rice go-embed
package pandoc

import "github.com/GeertJohan/go.rice"

type Resource struct {
	FileName string
	Content  string
}

var latexTemplateResources []Resource

func init() {
	templateBox := rice.MustFindBox("../_fixtures/latex_template")

	getResource := func(filename string) Resource {
		return Resource{
			FileName: filename,
			Content:  templateBox.MustString(filename),
		}
	}

	structureTex := getResource("structure.tex")
	cciconsSty := getResource("ccicons.sty")
	templateTex := getResource("template.tex")
	latexTemplateResources = []Resource{
		structureTex,
		cciconsSty,
		templateTex,
	}
}
