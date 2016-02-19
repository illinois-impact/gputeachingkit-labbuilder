package pkg

import (
	"github.com/k0kubun/pp"
	pf "gitlab.com/abduld/wgx-labpdf/pkg/pandocfilter"
	"log"
	"strconv"
	"strings"
)

var (
	LabTitle   = "Lab title..."
	LabModule  = -1
	LabAuthor  = "GPU Teaching Kit -- Accelerated Computing"
	LabNumber  = -1
	firstVisit = true
)

func LabNumberFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	if !firstVisit {
		return nil
	}
	firstVisit = false
	info, ok := meta.(map[string]interface{})
	if !ok {
		pp.Println(meta)
	}
	if title, ok := info["title"]; ok {
		LabTitle = pf.Stringify(title)
	} else {
		log.Fatal("Cannot find document title in title block.\n")
	}
	if module, ok := info["module"]; ok {
		mod := pf.Stringify(module)
		mod = strings.TrimSpace(mod)
		n, err := strconv.Atoi(mod)
		if err != nil {
			log.Fatal("The module field in the title is set to '" + mod + "'. Expecting a number.\n")
		}
		LabModule = n
	} else {
		log.Fatal("Cannot find module number in title block.")
	}
	if author, ok := info["author"]; ok {
		LabAuthor = pf.Stringify(author)
	} else {
		log.Println(
			"The author field in the title header is not set. Using the default author: '" +
				LabAuthor + "'.")
	}

	return nil
}

func init() {
	pf.Str("fds")
	AddFilter(LabNumberFilter)
}
