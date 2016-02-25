package build

import (
	"io/ioutil"

	"path/filepath"

	"github.com/cheggaaa/pb"
)

func Markdown(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
	doc, err := makeDoc(outputDir, cmakeFile, progress)
	if err != nil {
		return "", err
	}
	if progress == nil {
		progress = newProgressBar(doc.FileName)
		defer progress.Finish()
	}

	progress.Postfix("Creating the markdown file...")
	document, err := doc.markdown()
	if err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the markdown file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.Postfix("Writing Markdown file...")
	outFile := filepath.Join(outputDir, doc.FileName+".md")
	ioutil.WriteFile(outFile, []byte(document), 0644)
	incrementProgress(progress)

	return outFile, nil
}
