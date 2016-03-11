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

func Markdown(outputDir, cmakeFile string, progress *pb.ProgressBar) (string, error) {
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

	progressPostfix(progress, "Writing Markdown file...")
	outFile := filepath.Join(outputDir, "Module["+strconv.Itoa(doc.Module)+"]-"+doc.FileName+".md")

	tmpDir := os.TempDir()
	tmpOutFile := filepath.Join(tmpDir, "wgx-pandoc-markdown.md")
	ioutil.WriteFile(tmpOutFile, []byte(document), 0644)
	args := []string{
		"-s",
		"-N",
		"-f",
		pandoc.MarkdownFormat,
		"-t",
		pandoc.MarkdownFormat,
		"-o",
		outFile,
	}
	args = append(args, pandoc.DefaultFilter...)
	cmd := exec.Command("pandoc", args...)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		ioutil.WriteFile(filepath.Join(outputDir, doc.FileName+".gen.markdown.log"), out, 0644)
	}
	if err != nil {
		progress.FinishPrint("✖ Failed to generate TeX file. Error :: " + err.Error())
		return "", err
	}

	incrementProgress(progress)

	return outFile, nil
}
