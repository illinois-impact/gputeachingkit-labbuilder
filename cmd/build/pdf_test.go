package build

import (
	"path/filepath"
	"testing"
)

const techingkitDir = "~/gputeachingkit-labs-develop"

func TestBuildLab(t *testing.T) {
	cmakeFile := filepath.Join(techingkitDir, "/Module2/DeviceQuery/sources.cmake")
	PDF("/tmp", cmakeFile, nil)
}
func TestBuildAll(t *testing.T) {
	All("pdf", "/tmp", techingkitDir)
}
