package build

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/k0kubun/pp"
	"github.com/mattn/go-zglob"
	"github.com/mitchellh/go-homedir"
)

var (
	showProgress   bool
	filterDocument bool
)

func All(targetType string, outputDir0 string, showProgress0 bool, filterDocument0 bool, inputDir string) error {
	showProgress = showProgress0
	filterDocument = filterDocument0
	rootDir, _ := homedir.Expand(inputDir)
	matches, err := zglob.Glob(filepath.Join(rootDir, "**", "sources.cmake"))
	if err != nil {
		return err
	}
	outputDir, _ := homedir.Expand(outputDir0)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.WithField("output", outputDir).Error("Cannot create output directory")
		return err
	}

	if len(matches) == 0 {
		msg := "Invalid input directory. The input directory must contain one or more sources.cmake"
		log.Error(msg)
		return errors.New(msg)
	}

	var wg sync.WaitGroup
	wg.Add(len(matches))

	bars := []*pb.ProgressBar{}
	var barsMutex sync.Mutex
	for ii := range matches {
		cmakeFile := matches[ii]
		func() {
			defer wg.Done()
			if !path.IsAbs(cmakeFile) {
				cmakeFile = path.Join(rootDir, cmakeFile)
			}
			if !isCmakeLab(cmakeFile) {
				return
			}
			descriptionData, err := ioutil.ReadFile(
				filepath.Join(filepath.Dir(cmakeFile), "description.markdown"),
			)
			if err != nil {
				log.Panic("Cannot read config file for " + cmakeFile)
			}
			labName, err := getLabNameFromMarkdown(string(descriptionData))
			if err != nil {
				log.Panic("Cannot get lab name for " + cmakeFile)
			}
			bar := newProgressBar(labName)
			barsMutex.Lock()
			bars = append(bars, bar)
			barsMutex.Unlock()
			switch targetType {
			case "pdf":
				PDF(outputDir, cmakeFile, bar)
			case "markdown":
				Markdown(outputDir, cmakeFile, bar)
			case "html":
				HTML(outputDir, cmakeFile, bar)
			default:
				log.Panic("Does not understand how to make " + targetType + ". Valid target types are pdf, html, and markdown.")
			}

		}()
	}
	pool, err := pb.StartPool(bars...)

	wg.Wait()

	pool.Stop()
	return nil
}

func init() {
	log.SetLevel(log.DebugLevel)
	if false {
		pp.Println("dummy")
	}
}
