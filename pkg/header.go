package pkg


import (
	pf "github.com/oltolm/go-pandocfilters"
)

func HeaderFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	if k == "Header" {
		var level int = v.([]interface{})[0].(int)
		var attrs []interface{} = v.([]interface{})[1].([]interface{})
		var inlines []interface{} = v.([]interface{})[2].([]interface{})

		return pf.Header(level-1, attrs, inlines)
	}

	return nil
}

func init() {
	AddFilter(HeaderFilter)
}

