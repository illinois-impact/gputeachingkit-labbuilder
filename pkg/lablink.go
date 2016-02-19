package pkg

import (
	pf "gitlab.com/abduld/wgx-labpdf/pkg/pandocfilter"
	"strconv"
)

const BitbucketURL = "https://bitbucket.org/hwuligans/gputeachingkit-labs/src/master/"

func LabLinkFilter(k string, v interface{}, format string, meta interface{}) interface{} {

	if k == "Str" {
		value := v.(string)
		if value == "LINKTOLAB" {
			return pf.Str(BitbucketURL + strconv.Itoa(LabNumber))
		}
	}

	return nil
}

func init() {
	AddFilter(LabLinkFilter)
}
