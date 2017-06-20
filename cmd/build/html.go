package build

import (
	"io/ioutil"

	"path/filepath"

	"os"
	"os/exec"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/webgpu/gputeachingkit-labbuilder/pkg"
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
		progress.FinishPrint("✖ Failed " + doc.FileName + " to create the markdown file. Error :: " + err.Error())
		return "", err
	}
	incrementProgress(progress)

	progressPostfix(progress, "Writing Markdown file...")
	outFile := filepath.Join(outputDir, "Module["+strconv.Itoa(doc.Module)+"]-"+doc.FileName+".html")

	tmpDir := os.TempDir()
	tmpOutFile := filepath.Join(tmpDir, "gputeachingkit-labbuilder-markdown.md")
	ioutil.WriteFile(tmpOutFile, []byte(document), 0644)

	args := []string{
		"-s",
		"-N",
		"-f",
		pandoc.MarkdownFormat,
		"-c",
		cssResource.fileName,
		"--katex",
		"-t",
		"html5",
		"--template=" + htmlTemplate.fileName,
		"-o",
		outFile,
		tmpOutFile,
	}
	args = append(args, pandoc.DefaultFilter...)
	cmd := exec.Command("pandoc", args...)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		ioutil.WriteFile(filepath.Join(outputDir, doc.FileName+".gen.html.log"), out, 0644)
	}
	if err != nil {
		progress.FinishPrint("✖ Failed to generate HTML file. Error :: " + err.Error())
		return "", err
	}

	incrementProgress(progress)

	return outFile, nil
}
