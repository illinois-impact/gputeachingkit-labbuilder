package pkg

import (
	"strconv"

	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
)

const BitbucketURLEval = "https://bitbucket.org/hwuligans/gputeachingkit-labs/src/master/"
const BitbucketURLFull = "https://bitbucket.org/hwuligans/gputeachingkit-labs-full/src/master/"

func LabLinkFilter(k string, v interface{}, format string, meta interface{}) interface{} {

	if k == "Str" {
		value := v.(string)
		if value == "LINKTOLAB" {
			var url string
			if Config.IsFullTookit {
				url = BitbucketURLFull
			} else {
				url = BitbucketURLEval
			}
			return pf.Str(url + strconv.Itoa(Lab.Number))
		}
	}

	return nil
}

func init() {
	AddFilter(LabLinkFilter)
}
