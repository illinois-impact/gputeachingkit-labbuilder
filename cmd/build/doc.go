package build

import (
	"io/ioutil"

	"errors"
	"path/filepath"
	"regexp"

	"bytes"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/mitchellh/go-homedir"
)

func (d *doc) markdown() (string, error) {
	var document bytes.Buffer
	tmpl := template.Must(template.New(d.Name + "_template").Parse(markdownTemplate))
	err := tmpl.Execute(&document, d)
	if err != nil {
		return "", err
	}
	return document.String(), nil
}

func (d *doc) tex() (string, error) {
	var document bytes.Buffer
	tmpl := template.Must(template.New(d.Name + "_template").Parse(texTemplate))
	err := tmpl.Execute(&document, d)
	if err != nil {
		return "", err
	}
	return document.String(), nil
}

func makeDoc(outputDir, cmakeFile string, progress *pb.ProgressBar) (*doc, error) {
	if progress != nil {
		defer progress.Finish()
	}
	cmakeFile, _ = homedir.Expand(cmakeFile)
	if !isCmakeLab(cmakeFile) {
		return nil, errors.New("Invalid cmake file " + cmakeFile)
	}
	re := regexp.MustCompile(`add_lab\("(?P<lab_name>.*)"\)`)

	buf, err := ioutil.ReadFile(cmakeFile)
	if err != nil {
		if progress != nil {
			progress.FinishPrint("Failed to read cmake file")
		}
		return nil, err
	}
	content := string(buf)
	match := re.FindStringSubmatch(content)
	if len(match) != 2 {
		if progress != nil {
			progress.FinishPrint("Cannot parse the add_lab() line in " + cmakeFile)
		}
		return nil, err
	}
	fileName := match[1]
	if progress == nil {
		progress = newProgressBar(fileName)
		defer progress.Finish()
	}

	rootDir := filepath.Dir(cmakeFile)

	progress.Postfix("Starting ...")
	//configFileName := filepath.Join(rootDir, "config.json")
	descriptionFileName := filepath.Join(rootDir, "description.markdown")
	questionsFileName := filepath.Join(rootDir, "questions.json")
	answersFileName := filepath.Join(rootDir, "answers.json")
	codeTemplateFileName := filepath.Join(rootDir, "template.cu")
	codeSolutionFileName := filepath.Join(rootDir, "solution.cu")

	incrementProgress(progress)
	if !isFile(codeTemplateFileName) {
		codeTemplateFileName = filepath.Join(rootDir, "template.cpp")
	}
	if !isFile(codeSolutionFileName) {
		codeSolutionFileName = filepath.Join(rootDir, "solution.cpp")
	}

	progress.Postfix("Reading Files ...")
	readFile := func(pth string) string {
		if err != nil {
			return ""
		}
		if !isFile(pth) {
			err = errors.New("File " + pth + " not found")
			return ""
		}
		var buf []byte
		buf, err = ioutil.ReadFile(pth)
		if err != nil {
			return ""
		}
		return string(buf)
	}
	//config := readFile(configFileName)
	description := readFile(descriptionFileName)
	questions := readFile(questionsFileName)
	answers := readFile(answersFileName)
	codeTemplate := readFile(codeTemplateFileName)
	codeSolution := readFile(codeSolutionFileName)

	incrementProgress(progress)

	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while reading the files. Error :: " + err.Error())
		return nil, err
	}

	progress.Postfix("Getting lab name ...")
	labName, err := getLabNameFromMarkdown(string(description))
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while getting the lab name. Error :: " + err.Error())
		return nil, err
	}
	incrementProgress(progress)

	progress.Postfix("Removing title section ...")
	description = removeTitleYaml(description)
	incrementProgress(progress)

	progress.Postfix("Getting module number ...")
	moduleNumber, err := getModuleNumber(rootDir)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while getting the lab module number. Error :: " + err.Error())
		return nil, err
	}
	incrementProgress(progress)

	progress.Postfix("Getting questions and answers ...")
	questionAnswers, err := getQuestionsAnswers(questions, answers)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while getting the questions and answers. Error :: " + err.Error())
		return nil, err
	}
	incrementProgress(progress)

	return &doc{
		Module:          moduleNumber,
		FileName:        fileName,
		Name:            labName,
		Description:     description,
		QuestionAnswers: questionAnswers,
		CodeTemplate:    codeTemplate,
		CodeSolution:    codeSolution,
	}, nil
}

func init() {
	if false {
		log.Warn("dummy")
	}
}
