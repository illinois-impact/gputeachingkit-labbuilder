package build

import (
	"io/ioutil"
	"path/filepath"

	"github.com/cheggaaa/pb"
	"github.com/russross/blackfriday"
)

var (
	blackfridayExtensions int
	htmlFlags             int
)

func HTML(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
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

	progressPostfix(progress, "Building HTML file...")
	htmlRenderer := blackfriday.HtmlRenderer(htmlFlags, doc.Name, "")
	data := blackfriday.Markdown([]byte(document), htmlRenderer, blackfridayExtensions)
	if err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create pdf output. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progressPostfix(progress, "Copying the output file to destination directory...")
	outFile := filepath.Join(outputDir, doc.FileName+".html")
	if err = ioutil.WriteFile(outFile, data, 0644); err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to write the output file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.FinishPrint("✔ Completed " + doc.Name + " placing target at " + outFile)
	return outFile, nil

}

func init() {
	blackfridayExtensions = 0
	blackfridayExtensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	blackfridayExtensions |= blackfriday.EXTENSION_TABLES
	blackfridayExtensions |= blackfriday.EXTENSION_FENCED_CODE
	blackfridayExtensions |= blackfriday.EXTENSION_AUTOLINK
	blackfridayExtensions |= blackfriday.EXTENSION_STRIKETHROUGH
	blackfridayExtensions |= blackfriday.EXTENSION_SPACE_HEADERS

	htmlFlags = 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS

	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
}
