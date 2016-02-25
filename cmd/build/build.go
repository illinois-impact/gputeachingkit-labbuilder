package build

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/mattn/go-zglob"
	"github.com/mitchellh/go-homedir"
)

func All(targetType, outputDir, inputDir string) error {
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
}
