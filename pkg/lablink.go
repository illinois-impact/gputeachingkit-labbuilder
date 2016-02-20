package pkg

import (
	"strconv"

	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
)

const BitbucketURL = "https://bitbucket.org/hwuligans/gputeachingkit-labs/src/master/"

func LabLinkFilter(k string, v interface{}, format string, meta interface{}) interface{} {

	if k == "Str" {
		value := v.(string)
		if value == "LINKTOLAB" {
			return pf.Str(BitbucketURL + strconv.Itoa(Lab.Number))
		}
	}

	return nil
}

func init() {
	AddFilter(LabLinkFilter)
}
