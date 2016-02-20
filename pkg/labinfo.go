package pkg

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/k0kubun/pp"
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
	"golang.org/x/net/context"
)

type lab struct {
	Title  string
	Module int
	Author string
	Number int
}

var Lab = lab{}

func LabNumberFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	if _, ok := ctx.Value("VisitedLabNumberFilter").(bool); ok {
		return nil
	}
	ctx = context.WithValue(ctx, "VisitedLabNumberFilter", true)

	info, ok := meta.(map[string]interface{})
	if !ok {
		pp.Println(meta)
	}
	if title, ok := info["title"]; ok {
		Lab.Title = pf.Stringify(title)
	} else {
		logrus.Fatal("Cannot find document title in title block.\n")
	}
	if module, ok := info["module"]; ok {
		mod := pf.Stringify(module)
		mod = strings.TrimSpace(mod)
		n, err := strconv.Atoi(mod)
		if err != nil {
			logrus.Fatal("The module field in the title is set to '" + mod + "'. Expecting a number.\n")
		}
		Lab.Module = n
	} else {
		logrus.Fatal("Cannot find module number in title block.")
	}
	if author, ok := info["author"]; ok {
		Lab.Author = pf.Stringify(author)
	} else {
		logrus.Println(
			"The author field in the title header is not set. Using the default author: '" +
				Lab.Author + "'.")
	}

	return nil
}

func init() {
	//AddFilter(LabNumberFilter)
}
