package filter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"gitlab.com/abduld/wgx-pandoc/pkg"
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
	"gitlab.com/abduld/wgx-utils"
)

func toJSON(inputFilePath string) (string, error) {
	tmpDir := os.TempDir()
	outputFile := filepath.Join(tmpDir, "pandocJsonOutput.json")
	log.Debug("Generating pandoc json from markdown file")
	cmd := exec.Command("pandoc",
		"-o",
		outputFile,
		"-f",
		"markdown+hard_line_breaks+pandoc_title_block+lists_without_preceding_blankline+compact_definition_lists",
		"-t",
		"json",
		inputFilePath,
	)
	cmd.Dir = tmpDir
	buf, err := cmd.CombinedOutput()
	log.WithError(err).WithField("command_out", string(buf)).Debug("Ran pandoc to json command")
	if err != nil {
		return "", errors.New(string(buf) + " .. Error: " + err.Error())
	}
	return outputFile, err
}

func fromJSON(outputFilePath, inputFilePath string) error {
	tmpDir := os.TempDir()
	cmd := exec.Command("pandoc",
		"-o",
		outputFilePath,
		"-f",
		"json",
		"-t",
		"markdown+hard_line_breaks+pandoc_title_block+lists_without_preceding_blankline+compact_definition_lists",
		"-S",
		"-s",
		inputFilePath,
	)
	cmd.Dir = tmpDir
	buf, err := cmd.CombinedOutput()
	log.WithError(err).WithField("command_out", string(buf)).Debug("Ran pandoc to markdown command")
	if err != nil {
		return errors.New(string(buf) + " .. Error: " + err.Error())
	}
	return err
}

func fileName(pth string) string {
	base := filepath.Base(pth)
	return base[:len(base)-len(filepath.Ext(pth))]
}

func isMarkdownExt(pth string) bool {
	return filepath.Ext(pth) == ".markdown" ||
		filepath.Ext(pth) == ".md"
}

func Filter(outputFileDir, inputFilePath string, format string) (string, error) {

	var (
		doc               interface{}
		jsonInputFilePath string
	)

	inputFilePath, _ = homedir.Expand(inputFilePath)
	inputFilePath, _ = filepath.Abs(inputFilePath)
	jsonOutpuFilePath := filepath.Join(outputFileDir, fileName(inputFilePath)+"-filter.json")
	outputFilePath := filepath.Join(outputFileDir, fileName(inputFilePath)+"-filter.markdown")
	log.Debug("Input file is set to " + inputFilePath)
	log.Debug("Output file is set to " + outputFilePath)

	if !utils.IsFile(inputFilePath) {
		return "", errors.New("input file does not exist")
	}

	if isMarkdownExt(inputFilePath) {
		var err error
		log.Debug("File has a markdown extension... converting to JSON format.")
		jsonInputFilePath, err = toJSON(inputFilePath)
		if err != nil {
			return "", err
		}
	} else {
		jsonInputFilePath = inputFilePath
	}

	pandoc.Clear()

	inputData, err := ioutil.ReadFile(jsonInputFilePath)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(inputData, &doc); err != nil {
		return "", err
	}

	meta := doc.([]interface{})[0].(map[string]interface{})["unMeta"]
	for _, filter := range pandoc.Filters {
		doc = pf.Walk(doc, filter, format, meta)
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(jsonOutpuFilePath, b, 0644); err != nil {
		return "", err
	}

	if isMarkdownExt(inputFilePath) {
		if err := fromJSON(outputFilePath, jsonOutpuFilePath); err != nil {
			return "", err
		}
	}

	return outputFilePath, nil
}
