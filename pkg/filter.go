package pkg

import (
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
	"sync"
"golang.org/x/net/context"
	"github.com/Sirupsen/logrus"
)

var (
	Filters = []pf.Action{}
	mutex   sync.Mutex
	ctx context.Context
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
