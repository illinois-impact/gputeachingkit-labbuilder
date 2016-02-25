//go:generate rice embed-go
package build

import (
	"io/ioutil"
	"os"

	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/mattn/go-zglob"
	"github.com/mitchellh/go-homedir"
)

type questionAnswer struct {
	Question string
	Answer   string
}

type doc struct {
	Module          int
	FileName        string
	Name            string
	Description     string
	QuestionAnswers []questionAnswer
	CodeTemplate    string
	CodeSolution    string
}

type resource struct {
	fileName string
	content  string
}

var (
	box          = rice.MustFindBox("./../../_fixtures")
	templateData = box.MustString("tex.template")

	latexTemplateResources map[string]resource
)

func isFile(file string) bool {
	if fi, err := os.Stat(file); err != nil || fi.IsDir() {
		return false
	}
	return true
}

func copyFile(trgt, src string) error {
	// open files r and w
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(trgt)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err = io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func getLabName(configData string) (string, error) {
	var config struct {
		Name string `json:"string`
	}

	if err := json.Unmarshal([]byte(configData), &config); err != nil {
		return "", err
	}
	return config.Name, nil
}

func getQuestionsAnswers(questionsData, answeresData string) ([]questionAnswer, error) {
	var questions struct {
		Questions []string `json:"questions"`
	}
	var answers struct {
		Answers []string `json:"answers"`
	}
	qas := []questionAnswer{}

	if err := json.Unmarshal([]byte(questionsData), &questions); err != nil {
		return qas, err
	}
	if err := json.Unmarshal([]byte(answeresData), &answers); err != nil {
		return qas, err
	}
	for ii, question := range questions.Questions {
		item := questionAnswer{Question: question}
		if ii < len(answers.Answers) {
			item.Answer = answers.Answers[ii]
		}
		qas = append(qas, item)
	}
	return qas, nil
}

func getModuleNumber(path string) (int, error) {
	re := regexp.MustCompile(`Module(\d+)/`)
	match := re.FindStringSubmatch(path)
	if len(match) != 2 {
		return 0, errors.New("Cannot detect the module number from " + path)
	}
	return strconv.Atoi(match[1])
}

func removeTitleYaml(description string) string {
	var start, end int
	lines := strings.Split(description, "\n")
	for ii, line := range lines {
		if strings.HasPrefix(line, "---") {
			start = ii
			break
		}
	}
	for ii, line := range lines[start+1:] {
		if strings.HasPrefix(line, "---") {
			end = start + ii + 1
			break
		}
	}
	return strings.Join(lines[end+1:], "\n")
}

func writeLatexResources(dir string) {
	for _, res := range latexTemplateResources {
		ioutil.WriteFile(filepath.Join(dir, res.fileName), []byte(res.content), 0644)
	}
}

func buildPDF(info doc, document string, progress *pb.ProgressBar) (string, error) {
	progress.Postfix("Creating temporary directory...")
	tmpDir, err := ioutil.TempDir("", info.FileName+"-build")
	if err != nil {
		progress.FinishPrint("✖ Failed to create temporary directory. Error :: " + err.Error())
	}
	progress.Increment()

	//defer os.RemoveAll(dir) // clean up
	fileBaseName := filepath.Join(tmpDir, info.FileName)
	mdFileName := fileBaseName + ".md"
	texFileName := fileBaseName + ".tex"
	pdfFileName := fileBaseName + ".pdf"

	progress.Postfix("Writing resources to temporary directory...")
	writeLatexResources(tmpDir)
	ioutil.WriteFile(mdFileName, []byte(document), 0644)
	progress.Increment()

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
	progress.Increment()

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
	progress.Increment()

	return pdfFileName, nil
}

func newProgressBar(prefix string) *pb.ProgressBar {
	progress := pb.New(17)
	progress.Prefix(prefix)
	progress.SetWidth(80)
	progress.AlwaysUpdate = true
	progress.ShowFinalTime = true
	return progress
}
func incr(progress *pb.ProgressBar) {
	progress.Increment()
	progress.Update()
}

func isCmakeLab(cmakeFile string) bool {
	buf, err := ioutil.ReadFile(cmakeFile)
	if err != nil {
		return false
	}
	if !strings.Contains(string(buf), "add_lab(") {
		return false
	}
	return true
}
func Lab(outputDir, cmakeFile string, progress *pb.ProgressBar) {
	if progress != nil {
		defer progress.Finish()
	}
	cmakeFile, _ = homedir.Expand(cmakeFile)
	if !isCmakeLab(cmakeFile) {
		return
	}
	re := regexp.MustCompile(`add_lab\("(?P<lab_name>.*)"\)`)

	buf, err := ioutil.ReadFile(cmakeFile)
	if err != nil {
		if progress != nil {
			progress.FinishPrint("Failed to read cmake file")
		}
		return
	}
	content := string(buf)
	match := re.FindStringSubmatch(content)
	if len(match) != 2 {
		if progress != nil {
			progress.FinishPrint("Cannot parse the add_lab() line in " + cmakeFile)
		}
		return
	}
	fileName := match[1]
	if progress == nil {
		progress = newProgressBar(fileName)
		defer progress.Finish()
	}

	rootDir := filepath.Dir(cmakeFile)

	progress.Postfix("Starting ...")
	configFileName := filepath.Join(rootDir, "config.json")
	descriptionFileName := filepath.Join(rootDir, "description.markdown")
	questionsFileName := filepath.Join(rootDir, "questions.json")
	answersFileName := filepath.Join(rootDir, "answers.json")
	codeTemplateFileName := filepath.Join(rootDir, "template.cu")
	codeSolutionFileName := filepath.Join(rootDir, "solution.cu")

	progress.Increment()
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
	config := readFile(configFileName)
	description := readFile(descriptionFileName)
	questions := readFile(questionsFileName)
	answers := readFile(answersFileName)
	codeTemplate := readFile(codeTemplateFileName)
	codeSolution := readFile(codeSolutionFileName)

	progress.Increment()

	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while reading the files. Error :: " + err.Error())
		return
	}

	progress.Postfix("Removing title section ...")
	description = removeTitleYaml(description)
	progress.Increment()

	progress.Postfix("Getting module number ...")
	moduleNumber, err := getModuleNumber(rootDir)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while getting the lab module number. Error :: " + err.Error())
		return
	}
	progress.Increment()

	progress.Postfix("Getting lab name ...")
	labName, err := getLabName(config)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while getting the lab name. Error :: " + err.Error())
		return
	}
	progress.Increment()

	progress.Postfix("Getting questions and answers ...")
	questionAnswers, err := getQuestionsAnswers(questions, answers)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " while getting the questions and answers. Error :: " + err.Error())
		return
	}
	progress.Increment()

	docVars := doc{
		Module:          moduleNumber,
		FileName:        fileName,
		Name:            labName,
		Description:     description,
		QuestionAnswers: questionAnswers,
		CodeTemplate:    codeTemplate,
		CodeSolution:    codeSolution,
	}

	var document bytes.Buffer

	progress.Postfix("Creating the markdown file...")
	tmpl := template.Must(template.New(labName + "_template").Parse(templateData))
	err = tmpl.Execute(&document, docVars)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " to create the markdown file. Error :: " + err.Error())
		return
	}
	progress.Increment()

	progress.Postfix("Building PDF file...")
	pdfFile, err := buildPDF(docVars, document.String(), progress)
	if err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " to create pdf output. Error :: " + err.Error())
		return
	}
	progress.Increment()

	progress.Postfix("Copying the output file to destination directory...")
	outFile := filepath.Join(outputDir, docVars.FileName+".pdf")
	if err = copyFile(outFile, pdfFile); err != nil {
		progress.FinishPrint("✖ Failed " + fileName + " to copy the output file. Error :: " + err.Error())
		return
	}
	incr(progress)

	progress.FinishPrint("✔ Completed " + labName + " placing target at " + outFile)

}

func All(outputDir, inputDir string) error {
	rootDir, _ := homedir.Expand(inputDir)
	matches, err := zglob.Glob(filepath.Join(rootDir, "**", "sources.cmake"))
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(matches))

	bars := []*pb.ProgressBar{}
	var barsMutex sync.Mutex
	for ii := range matches {
		cmakeFile := filepath.Join(rootDir, matches[ii])
		go func() {
			defer wg.Done()
			if !isCmakeLab(cmakeFile) {
				return
			}
			configData, err := ioutil.ReadFile(filepath.Join(filepath.Dir(cmakeFile), "config.json"))
			if err != nil {
				log.Panic("Cannot read config file for " + cmakeFile)
			}
			labName, err := getLabName(string(configData))
			if err != nil {
				log.Panic("Cannot get lab name for " + cmakeFile)
			}

			bar := newProgressBar(labName)
			barsMutex.Lock()
			bars = append(bars, bar)
			barsMutex.Unlock()

			Lab(outputDir, cmakeFile, bar)

		}()
	}
	pool, err := pb.StartPool(bars...)

	wg.Wait()

	pool.Stop()
	return nil
}

func init() {
	log.SetLevel(log.DebugLevel)

	latexTemplateBox := rice.MustFindBox("./../../_fixtures/latex_template")
	getTexResource := func(filename string) resource {
		return resource{
			fileName: filename,
			content:  latexTemplateBox.MustString(filename),
		}
	}
	latexTemplateResources = map[string]resource{
		"structure.tex": getTexResource("structure.tex"),
		"ccicons.sty":   getTexResource("ccicons.sty"),
		"template.tex":  getTexResource("template.tex"),
	}
}
