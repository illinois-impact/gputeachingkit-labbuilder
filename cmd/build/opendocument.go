package build

import (
	"io/ioutil"

	"path/filepath"

	"os"
	"os/exec"
	"strconv"

	"github.com/cheggaaa/pb"
	"gitlab.com/abduld/wgx-pandoc/pkg"
)

const docxCopyright = `
------

This work is licensed by UIUC and NVIDIA (2016) under a Creative Commons Attribution-NonCommercial 4.0 License.
`

func OpenDocument(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {

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
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the markdown file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progressPostfix(progress, "Writing ODT file...")
	outFile := filepath.Join(outputDir, "Module["+strconv.Itoa(doc.Module)+"]-"+doc.FileName+".odt")

	tmpDir := os.TempDir()
	tmpOutFile := filepath.Join(tmpDir, "wgx-pandoc-markdown.markdown")
	ioutil.WriteFile(tmpOutFile, []byte(document), 0644)

	args := []string{
		"-s",
		"-N",
		"-f",
		pandoc.MarkdownFormat,
		"-t",
		"opendocument",
		"--template=" + opendocumentTemplate.fileName,
		"-o",
		outFile,
		tmpOutFile,
	}
	args = append(args, pandoc.DefaultFilter...)
	cmd := exec.Command("pandoc", args...)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		ioutil.WriteFile(filepath.Join(outputDir, doc.FileName+".gen.opendocument.log"), out, 0644)
	}
	if err != nil {
		progress.FinishPrint("✖ Failed to generate OpenDocument file. Error :: " + err.Error())
		return "", err
	}

	incrementProgress(progress)
    
	return outFile, nil
}

func Docx(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {

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
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the markdown file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)
    

	progressPostfix(progress, "Writing Docx file...")
    
	tmpDir := os.TempDir()
	tmpOutFile := filepath.Join(tmpDir, "wgx-pandoc-docx-markdown.markdown")
	ioutil.WriteFile(tmpOutFile, []byte(document+docxCopyright), 0644)
	outFile := filepath.Join(outputDir, "Module["+strconv.Itoa(doc.Module)+"]-"+doc.FileName+".docx")

	args := []string{
		"-s",
		"-N",
		"-f",
		pandoc.MarkdownFormat,
		"-t",
		"docx",
		"--template=" + opendocumentTemplate.fileName,
		"-o",
		outFile,
		tmpOutFile,
	}
	args = append(args, pandoc.DefaultFilter...)
	cmd := exec.Command("pandoc", args...)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		ioutil.WriteFile(filepath.Join(outputDir, doc.FileName+".gen.docx.log"), out, 0644)
	}
	if err != nil {
		progress.FinishPrint("✖ Failed to generate Docx file. Error :: " + err.Error())
		return "", err
	}

	return outFile, nil
}
