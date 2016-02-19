package main
import (
	pf "gitlab.com/abduld/wgx-labpdf/pkg/pandocfilter"
	"gitlab.com/abduld/wgx-labpdf/pkg"
)
func main() {
	pf.ToJsonFilter(pkg.HeaderFilter)
}
