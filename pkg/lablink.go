package pandoc

import (
	"strconv"

	pf "github.com/webgpu/gputeachingkit-labbuilder/pkg/pandocfilter"
)

func LabLinkFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	if k == "Link" {
		value := v.([]interface{})
		targets := value[2].([]interface{})
		target := targets[0].(string)
		if target == "LINKTOLAB" {
			var url string
			if Config.IsFullTookit {
				url = BitbucketURLFull
			} else {
				url = BitbucketURLEval
			}
			url = url + "Module" + strconv.Itoa(Lab.Module)
			targets[0] = url
			res := pf.Link(
				value[0].([]interface{}),
				value[1].([]interface{}),
				targets,
			)
			return res
		}
	}
	return nil
}

func init() {
	AddFilter(LabLinkFilter)
}
