package build

import (
	"io/ioutil"
	"path/filepath"

	"strconv"

	"github.com/cheggaaa/pb"
)

func marksy(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
	doc, err := makeDoc(outputDir, cmakeFile, progress)
	if err != nil {
		return "", err
	}
	if progress == nil {
		progress = newProgressBar(doc.FileName)
		defer progress.Finish()
	}

	progressPostfix(progress, "Creating the markdown file...")
	data, err := doc.markdown()
	if err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the tex file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	ext := ".marksy"
	outFile := filepath.Join(outputDir, "Module["+strconv.Itoa(doc.Module)+"]-"+doc.FileName+ext)
	if err = ioutil.WriteFile(outFile, []byte(data), 0644); err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to write the output file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.FinishPrint("✔ Completed " + doc.Name + " placing target at " + outFile)
	return outFile, nil

}

func MarksyText(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
	return marksy(outputDir, cmakeFile, progress)
}

func init() {
}
