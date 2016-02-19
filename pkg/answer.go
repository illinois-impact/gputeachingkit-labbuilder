package pkg

import "regexp"

var inanswer = false

func AnswerFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	pattern1 := regexp.MustCompile("<!-- BEGIN ANSWER -->")
	pattern2 := regexp.MustCompile("<!-- END ANSWER -->")
	if k == "RawBlock" {
		var format string = v.([]interface{})[0].(string)
		var s string = v.([]interface{})[1].(string)

		if format == "html" {
			if "" != pattern1.FindString(s) {
				inanswer = true
				return nil
			} else if "" != pattern2.FindString(s) {
				inanswer = false
				return nil
			}
		}
	}
	if inanswer {
		return []interface{}{}
	}
	return nil
}

func init() {
	AddFilter(AnswerFilter)
}
