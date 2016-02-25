package build

import (
	"path/filepath"
	"testing"
)

const techingkitDir = "~/gputeachingkit-labs-develop"

func TestBuildLab(t *testing.T) {
	Lab("/tmp", filepath.Join(techingkitDir, "/Module2/DeviceQuery/sources.cmake"), nil)
}
func TestBuildAll(t *testing.T) {
	All("/tmp", techingkitDir)
}
