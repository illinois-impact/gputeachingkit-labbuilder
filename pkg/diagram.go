package pandoc

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/bugsnag/osext"
	"github.com/k0kubun/pp"
	"github.com/spf13/cast"
	pf "github.com/webgpu/gputeachingkit-labbuilder/pkg/pandocfilter"
)

var (
	once         sync.Once
	mermaidExe   string
	phantomjsExe string
)

// see https://github.com/raghur/mermaid-filter

func DiagramFilter(k string, v interface{}, format string, meta interface{}) interface{} {

	if k != "CodeBlock" {
		return nil
	}
	once.Do(func() {
		findExe := func(dirName, exe string) (string, error) {
			exeDir, err := osext.ExecutableFolder()
			if err != nil {
				logrus.WithError(err).Warn("Cannot get executable folder")
			}
			if err == nil && runtime.GOOS != "windows" {
				var nodeModulesDir string
				if dir := filepath.Join(exeDir, "node_modules"); com.IsDir(dir) {
					nodeModulesDir = dir
				} else if dir := filepath.Join(exeDir, "..", "node_modules"); com.IsDir(dir) {
					nodeModulesDir = dir
				}

				if nodeModulesDir != "" && com.IsDir(filepath.Join(nodeModulesDir, dirName)) {
					return filepath.Join(nodeModulesDir, dirName, "bin", exe), nil
				}
			}
			return exec.LookPath(exe)
		}
		var err error
		if mermaidExe, err = findExe("mermaid", "mermaid.js"); err != nil {
			logrus.Warn("Mermaid is not installed. Disabling the diagram filter.")
			return
		}
		if phantomjsExe, err = findExe("phantomjs", "phantomjs"); err != nil {
			logrus.Warn("PhantomJS is not installed. Disabling the diagram filter.")
			return
		}
	})

	if mermaidExe == "" || phantomjsExe == "" {
		return nil
	}

	hasMemaidClass := false
	attrs := v.([]interface{})[0].([]interface{})
	content := v.([]interface{})[1].(string)
	classes := cast.ToStringSlice(attrs[1].([]interface{}))
	keyvals := [][]interface{}{}
	if val, ok := attrs[2].([][]interface{}); ok {
		keyvals = val
	}
	pp.Println(keyvals)
	for _, cls := range classes {
		if cls == "mermaid" {
			hasMemaidClass = true
		}
	}
	if !hasMemaidClass {
		return nil
	}
	opts := map[string]interface{}{
		"width":  500,
		"format": "png",
	}
	for _, keyval := range keyvals {
		key, ok := keyval[0].(string)
		if !ok {
			logrus.WithField("key", key).Fatal("Invalid key value in digram filter.")
			continue
		}
		if len(keyval) == 1 {
			opts[key] = true
		} else {
			opts[key] = keyval[1]
		}
	}

	tmpDir, err := ioutil.TempDir("", "pandocfilter")
	if err != nil {
		logrus.WithError(err).Fatal("Cannot create temporary directory in digram filter.")
	}
	//defer os.RemoveAll(dir) // clean up
	fileBaseName := filepath.Join(tmpDir, "mermaid")
	fileName := fileBaseName + ".mmd"

	if err = ioutil.WriteFile(fileName, []byte(content), 0644); err != nil {
		logrus.WithError(err).
			WithField("file", fileName).
			WithField("content", content).
			Fatal("Cannot write content to temporary file in digram filter.")
	}

	formatFlag := "-p"
	if val, ok := opts["format"]; ok {
		if val, ok := val.(string); ok {
			switch val {
			case "svg":
				formatFlag = "-s"
			case "png":
				formatFlag = "-p"
			default:
				logrus.WithField("format", val).Fatal("Invalid output format in digram filter.")
			}
		}
	}
	outputFile := fileBaseName
	switch formatFlag {
	case "-p":
		outputFile += ".png"
	case "-s":
		outputFile += ".svg"
	}

	cmd := exec.Command(mermaidExe, "-v", "-e", phantomjsExe, "-o", tmpDir, "-w", cast.ToString(opts["width"]), format, fileName)
	cmd.Dir = tmpDir
	cmd.Env = nil
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithError(err).WithField("command", strings.Join(cmd.Args, " ")).Fatal("Failed to generate output file.")
	}
	logrus.WithField("output_file", outputFile).
		WithField("command", strings.Join(cmd.Args, " ")).
		Debug("Command output = " + string(cmdOut))

	return pf.Para([]interface{}{
		pf.Image(
			[]interface{}{},
			[]interface{}{},
			[]interface{}{outputFile, ""},
		),
	})
}

func init() {
	AddFilter(DiagramFilter)
}
