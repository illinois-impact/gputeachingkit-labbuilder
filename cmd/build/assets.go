//go:generate rice embed-go

package build

import (
	"io/ioutil"
	"path/filepath"

	"github.com/GeertJohan/go.rice"
)

var (
	box                     = rice.MustFindBox("./../../_fixtures")
	texTemplate             = box.MustString("md.template")
	markdownTexTemplate     = box.MustString("tex.template")
	markdownRegularTemplate = box.MustString("md.template")

	latexTemplateResources map[string]resource
	htmlTemplate           resource
	rtfTemplate            resource
	opendocumentTemplate   resource
	cssResource            resource
)

func init() {
	tmpDir, _ := ioutil.TempDir("", "wgx-pandoc-assets")

	latexTemplateBox := rice.MustFindBox("./../../_fixtures/templates/latex")
	htmlTemplateBox := rice.MustFindBox("./../../_fixtures/templates/html")
	rtfTemplateBox := rice.MustFindBox("./../../_fixtures/templates/rtf")
	opendocumentTemplateBox := rice.MustFindBox("./../../_fixtures/templates/opendocument")
	getResource := func(box *rice.Box, baseName string) resource {
		data := box.MustString(baseName)
		fileName := filepath.Join(tmpDir, baseName)
		ioutil.WriteFile(fileName, box.MustBytes(baseName), 0644)
		return resource{
			fileName: fileName,
			baseName: baseName,
			content:  data,
		}
	}
	getTexResource := func(filename string) resource {
		return getResource(latexTemplateBox, filename)
	}
	latexTemplateResources = map[string]resource{
		"structure.tex": getTexResource("structure.tex"),
		"ccicons.sty":   getTexResource("ccicons.sty"),
		"template.tex":  getTexResource("template.tex"),
	}
	htmlTemplate = getResource(htmlTemplateBox, "template.html")
	cssResource = getResource(htmlTemplateBox, "style.css")
	rtfTemplate = getResource(rtfTemplateBox, "template.rtf")
	opendocumentTemplate = getResource(opendocumentTemplateBox, "template.opendocument")
}
