package build

import (
	"io/ioutil"
	"path/filepath"

	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/github_flavored_markdown"
)

var (
	blackfridayExtensions int
	htmlFlags             int
)

func blackfridayCommon(renderer blackfriday.Renderer, outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
	doc, err := makeDoc(outputDir, cmakeFile, progress)
	if err != nil {
		return "", err
	}
	if progress == nil {
		progress = newProgressBar(doc.FileName)
		defer progress.Finish()
	}

	progressPostfix(progress, "Creating the markdown file...")
	document, err := doc.markdown()
	if err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the tex file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	document = removeFrontMatter(document)

	var data []byte
	if _, ok := renderer.(*blackfriday.Latex); ok {
		progressPostfix(progress, "Building TeX file...")
		data = blackfriday.Markdown([]byte(document), renderer, blackfridayExtensions)
	} else {
		progressPostfix(progress, "Building HTML file...")
		if true {
			data = github_flavored_markdown.Markdown([]byte(document))
		} else {
			data = blackfriday.Markdown([]byte(document), renderer, blackfridayExtensions)
		}
	}
	if err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create pdf output. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progressPostfix(progress, "Copying the output file to destination directory...")

	ext := ".html"
	if _, ok := renderer.(*blackfriday.Latex); ok {
		ext = ".tex"
	}

	outFile := filepath.Join(outputDir, "Module["+strconv.Itoa(doc.Module)+"]-"+doc.FileName+ext)
	if err = ioutil.WriteFile(outFile, data, 0644); err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to write the output file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.FinishPrint("✔ Completed " + doc.Name + " placing target at " + outFile)
	return outFile, nil

}

func BlackfridayHTML(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
	renderer := blackfriday.HtmlRenderer(htmlFlags, "HTML", "")
	return blackfridayCommon(renderer, outputDir, cmakeFile, progress)
}

func BlackfridayTex(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
	renderer := blackfriday.LatexRenderer(0)
	return blackfridayCommon(renderer, outputDir, cmakeFile, progress)
}

func init() {
	blackfridayExtensions = 0
	blackfridayExtensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	blackfridayExtensions |= blackfriday.EXTENSION_TABLES
	blackfridayExtensions |= blackfriday.EXTENSION_FENCED_CODE
	blackfridayExtensions |= blackfriday.EXTENSION_AUTOLINK
	blackfridayExtensions |= blackfriday.EXTENSION_STRIKETHROUGH
	blackfridayExtensions |= blackfriday.EXTENSION_SPACE_HEADERS
	blackfridayExtensions |= blackfriday.EXTENSION_AUTO_HEADER_IDS

	htmlFlags = 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS

	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
}
