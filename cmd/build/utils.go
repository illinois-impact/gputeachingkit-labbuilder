package build

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"gitlab.com/abduld/wgx-md-frontmatter"
)

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

func getLabNameFromConfigJSON(configData string) (string, error) {
	var config struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal([]byte(configData), &config); err != nil {
		return "", err
	}
	return config.Name, nil
}

func getLabNameFromMarkdown(mk string) (string, error) {
	var front struct {
		Title string `yaml:"title"`
	}
	frontmatter := getFrontMatter(mk)
	if err := yaml.Unmarshal([]byte(frontmatter), &front); err != nil {
		log.Panic("Cannot get lab title")
		return "", err
	}
	return front.Title, nil
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

func getFrontMatter(description string) string {
	return frontmatter.Extract(description)
}

func removeFrontMatter(description string) string {
	return frontmatter.Trim(description)
}

func writeLatexResources(dir string) {
	for _, res := range latexTemplateResources {
		ioutil.WriteFile(filepath.Join(dir, res.fileName), []byte(res.content), 0644)
	}
}

func newProgressBar(prefix string) *pb.ProgressBar {

	progress := pb.New(17)
	progress.Prefix(prefix)
	progress.SetWidth(80)
	progress.Start()
	progress.SetRefreshRate(100 * time.Millisecond)
	progress.AlwaysUpdate = true
	progress.ShowFinalTime = true
	return progress
}
func incrementProgress(progress *pb.ProgressBar) {
	progress.Increment()
	progress.Update()
}

func progressPostfix(progress *pb.ProgressBar, s string) {
	progress.Increment()
	if showProgress {
		progress.Postfix(s)
	}
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
