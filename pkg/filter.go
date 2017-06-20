package pandoc

import (
	"sync"

	"github.com/sirupsen/logrus"
	pf "github.com/webgpu/gputeachingkit-labbuilder/pkg/pandocfilter"
	"golang.org/x/net/context"
)

var (
	Filters = []pf.Action{
		LabInfoFilter,
	}
	mutex sync.Mutex
	ctx   context.Context
)

func AddFilter(filter pf.Action) {
	mutex.Lock()
	defer mutex.Unlock()

	Filters = append(Filters, filter)
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	ctx = context.Background()
}

func Clear() {
	ctx = context.Background()
	Lab = lab{}
}
