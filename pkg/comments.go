//Pandoc filter that causes everything between
//'<!-- BEGIN COMMENT -->' and '<!-- END COMMENT -->'
//to be ignored.  The comment lines must appear on
//lines by themselves, with blank lines surrounding
//them.

package pkg

import "regexp"

var incomment = false

func CommentFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	pattern1 := regexp.MustCompile("<!-- BEGIN COMMENT -->")
	pattern2 := regexp.MustCompile("<!-- END COMMENT -->")
	if k == "RawBlock" {
		var format string = v.([]interface{})[0].(string)
		var s string = v.([]interface{})[1].(string)

		if format == "html" {
			if "" != pattern1.FindString(s) {
				incomment = true
				return nil
			} else if "" != pattern2.FindString(s) {
				incomment = false
				return nil
			}
		}
	}
	if incomment {
		return []interface{}{}
	}
	return nil
}

func init() {
	AddFilter(CommentFilter)
}
