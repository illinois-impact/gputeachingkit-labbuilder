// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/janeczku/go-spinner"
	"github.com/mattn/go-zglob"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [./path/to/GPUTeachingKit-Labs]",
	Short: "Makes the lab using the same mechanism as the make_lab_handout.py",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must provide a path to the base directory")
		}
		return nil
	},
	RunE: build,
}

func isFile(file string) bool {
	if fi, err := os.Stat(file); err != nil || fi.IsDir() {
		return false
	}
	return true
}

func buildLab(cmakeFile string) {
	buf, err := ioutil.ReadFile(cmakeFile)
	if err != nil {
		log.Debug("Cannot read " + cmakeFile)
		return
	}
	content := string(buf)
	if !strings.Contains(content, "add_lab(") {
		return
	}
	re := regexp.MustCompile(`add_lab\("(?P<lab_name>.*)"\)`)

	match := re.FindStringSubmatch(content)
	if len(match) != 2 {
		log.Error("Cannot parse the add_lab() line in " + cmakeFile)
		return
	}
	labName := match[1]
	progress := spinner.StartNew("Processing " + labName)
	progress.SetCharset([]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"})
	defer progress.Stop()

	baseFileName := filepath.Base(cmakeFile)

	configFileName := filepath.Join(baseFileName, "config.json")
	descriptionFileName := filepath.Join(baseFileName, "description.markdown")
	questionsFileName := filepath.Join(baseFileName, "questions.json")
	answersFileName := filepath.Join(baseFileName, "answers.json")
	codeTemplateFileName := filepath.Join(baseFileName, "template.cu")
	codeSolutionFileName := filepath.Join(baseFileName, "solution.cu")

	if !isFile(codeTemplateFileName) {
		codeTemplateFileName = filepath.Join(baseFileName, "template.cpp")
	}
	if !isFile(codeSolutionFileName) {
		codeSolutionFileName = filepath.Join(baseFileName, "solution.cpp")
	}

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
	if err != nil {
		log.Error(err)
		fmt.Println("✖ Failed " + labName)
	}
	fmt.Println("✔ Completed " + labName)
	//log.Println(content)
}
func build(cmd *cobra.Command, args []string) error {
	rootDir := args[0]
	matches, err := zglob.Glob(filepath.Join(rootDir, "**", "sources.cmake"))
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(matches))
	for _, match := range matches {
		go func(cmakeFile string) {
			defer wg.Done()
			buildLab(filepath.Join(rootDir, cmakeFile))
		}(match)
	}
	wg.Wait()
	return nil
}

func init() {
	log.SetLevel(log.DebugLevel)
	RootCmd.AddCommand(buildCmd)
}
