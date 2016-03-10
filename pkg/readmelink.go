package pandoc

import (
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
)

func ReadmeLinkFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	if k == "Link" {
		value := v.([]interface{})
		targets := value[2].([]interface{})
		target := targets[0].(string)
		if target == "LINKTOREADME" {
			var url string
			if Config.IsFullTookit {
				url = BitbucketURLFull
			} else {
				url = BitbucketURLEval
			}
			targets[0] = url + "README"
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
	AddFilter(ReadmeLinkFilter)
}
