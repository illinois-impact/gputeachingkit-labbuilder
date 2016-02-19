package pkg

import (
	pf "gitlab.com/abduld/wgx-labpdf/pkg/pandocfilter"
	"sync"
)

var (
	Filters = []pf.Action{}
	mutex   sync.Mutex
)

func AddFilter(filter pf.Action) {
	mutex.Lock()
	defer mutex.Unlock()

	Filters = append(Filters, filter)
}

func init() {

}
