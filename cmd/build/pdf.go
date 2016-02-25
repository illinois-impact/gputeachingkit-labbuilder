//go:generate rice embed-go
package build

import (
	"io/ioutil"

	"os/exec"
	"path/filepath"

	"github.com/cheggaaa/pb"
)

func buildPDF(doc *doc, document string, progress *pb.ProgressBar) (string, error) {
	progress.Postfix("Creating temporary directory...")
	tmpDir, err := ioutil.TempDir("", doc.FileName+"-build")
	if err != nil {
		progress.FinishPrint("✖ Failed to create temporary directory. Error :: " + err.Error())
	}
	incrementProgress(progress)

	//defer os.RemoveAll(dir) // clean up
	fileBaseName := filepath.Join(tmpDir, doc.FileName)
	mdFileName := fileBaseName + ".md"
	texFileName := fileBaseName + ".tex"
	pdfFileName := fileBaseName + ".pdf"

	progress.Postfix("Writing resources to temporary directory...")
	writeLatexResources(tmpDir)
	ioutil.WriteFile(mdFileName, []byte(document), 0644)
	incrementProgress(progress)

	progress.Postfix("Generating TeX file...")
	cmd := exec.Command("pandoc",
		"-s",
		"-N",
		"--template="+latexTemplateResources["template.tex"].fileName,
		mdFileName,
		"-o",
		texFileName,
	)
	cmd.Dir = tmpDir

	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		ioutil.WriteFile(fileBaseName+".gen.tex.log", out, 0644)
	}
	if err != nil {
		progress.FinishPrint("✖ Failed to generate TeX file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.Postfix("Generating PDF file...")
	cmd = exec.Command("pdflatex",
		texFileName,
		"-o",
		pdfFileName,
	)
	cmd.Dir = tmpDir

	out, err = cmd.CombinedOutput()
	if len(out) > 0 {
		ioutil.WriteFile(fileBaseName+".gen.pdf.log", []byte(out), 0644)
	}
	if err != nil {
		progress.FinishPrint("✖ Failed to generate PDF file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	return pdfFileName, nil
}

func PDF(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
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
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the tex file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.Postfix("Building PDF file...")
	pdfFile, err := buildPDF(doc, document, progress)
	if err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create pdf output. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.Postfix("Copying the output file to destination directory...")
	outFile := filepath.Join(outputDir, doc.FileName+".pdf")
	if err = copyFile(outFile, pdfFile); err != nil {
		progress.FinishPrint("✖ Failed " + doc.FileName + " to copy the output file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progress.FinishPrint("✔ Completed " + doc.Name + " placing target at " + outFile)
	return outFile, nil

}
