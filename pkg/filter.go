package pandoc

import (
	"sync"

	"github.com/Sirupsen/logrus"
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
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
