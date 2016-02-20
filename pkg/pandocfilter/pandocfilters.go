//Author: Oleg Tolmatcev <oleg.tolmatcev@gmail.com>
//Copyright: (C) 2013-2016 Oleg Tolmatcev
//
//Author: John MacFarlane <jgm@berkeley.edu>
//Copyright: (C) 2013 John MacFarlane
//License: BSD3
//
//Functions to aid writing python scripts that process the pandoc
//AST serialized as JSON.
package pandocfilters

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Action func(string, interface{}, string, interface{}) interface{}

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
}

//Walk a tree, applying an action to every object.
//Returns a modified tree.
func Walk(x interface{}, action Action, format string, meta interface{}) interface{} {
	switch x := x.(type) {
	case []interface{}:
		array := []interface{}{}
		for _, item := range x {
			if item1, ok := item.(map[string]interface{}); ok && item1["t"] != nil {
				res := action(item1["t"].(string), item1["c"], format, meta)
				if res == nil {
					array = append(array, Walk(item1, action, format, meta))
				} else if res1, ok := res.([]interface{}); ok {
					for _, z := range res1 {
						array = append(array, Walk(z, action, format, meta))
					}
				} else {
					array = append(array, Walk(res, action, format, meta))
				}
			} else {
				array = append(array, Walk(item, action, format, meta))
			}
		}
		return array
	case map[string]interface{}:
		obj := map[string]interface{}{}
		for k, _ := range x {
			obj[k] = Walk(x[k], action, format, meta)
		}
		return obj
	default:
		return x
	}
}

//Converts an action into a filter that reads a JSON-formatted
//pandoc document from stdin, transforms it by walking the tree
//with the action, and returns a new JSON-formatted pandoc document
//to stdout.  The argument is a function action(key, value, format, meta),
//where key is the type of the pandoc object (e.g. 'Str', 'Para'),
//value is the contents of the object (e.g. a string for 'Str',
//a list of inline elements for 'Para'), format is the target
//output format (which will be taken for the first command line
//argument if present), and meta is the document's metadata.
//If the function returns None, the object to which it applies
//will remain unchanged.  If it returns an object, the object will
//be replaced.  If it returns a list, the list will be spliced in to
//the list to which the target object belongs.  (So, returning an
//empty list deletes the object.)
func ToJsonFilter(action Action) {
	decoder := json.NewDecoder(os.Stdin)
	var doc []interface{}
	decoder.Decode(&doc)
	var format string
	if len(os.Args) > 1 {
		format = os.Args[1]
	} else {
		format = ""
	}
	altered := Walk(doc, action, format, doc[0].(map[string]interface{})["unMeta"])
	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(altered)
}

//Walks the tree x and returns concatenated string content,
//leaving out all formatting.
func Stringify(x interface{}) string {
	result := []string{}
	goFunc := func(key string, val interface{}, format string, meta interface{}) interface{} {
		if key == "Str" {
			result = append(result, val.(string))
		} else if key == "Code" {
			result = append(result, val.([]interface{})[1].(string))
		} else if key == "Math" {
			result = append(result, val.([]interface{})[1].(string))
		} else if key == "LineBreak" {
			result = append(result, " ")
		} else if key == "Space" {
			result = append(result, " ")
		}
		return nil
	}
	Walk(x, goFunc, "", nil)
	return strings.Join(result, "")
}

//Returns an attribute list, constructed from the
//dictionary attrs.
func Attributes(attrs map[string]interface{}) []interface{} {
	if attrs == nil {
		attrs = map[string]interface{}{}
	}
	ident, _ := attrs["id"].(string)
	classes, _ := attrs["classes"].([]string)
	if classes == nil {
		classes = []string{}
	}
	keyvals := [][]string{}
	for x, v := range attrs {
		if x != "classes" && x != "id" {
			keyvals = append(keyvals, []string{x, v.(string)})
		}
	}
	return []interface{}{ident, classes, keyvals}
}

// Constructors for block elements

func Plain(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Plain", "c": inlines}
}

func Para(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Para", "c": inlines}
}

func CodeBlock(attrs []interface{}, str string) map[string]interface{} {
	return map[string]interface{}{"t": "CodeBlock", "c": []interface{}{attrs, str}}
}

func RawBlock(format interface{}, str string) map[string]interface{} {
	return map[string]interface{}{"t": "RawBlock", "c": []interface{}{format, str}}
}

func BlockQuote(blocks []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "BlockQuote", "c": blocks}
}

func OrderedList(attrs []interface{}, listOfItems [][]interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "OrderedList", "c": []interface{}{attrs, listOfItems}}
}

func BulletList(listOfItems [][]interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "BulletList", "c": listOfItems}
}

func DefinitionList(defList []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "DefinitionList", "c": defList}
}

func Header(level int, attrs []interface{}, inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Header", "c": []interface{}{level, attrs, inlines}}
}

func HorizontalRule() map[string]interface{} {
	return map[string]interface{}{"t": "HorizontalRule", "c": []interface{}{}}
}

func Table(caption []interface{}, colAlign []interface{}, relColWidth []interface{},
	colHeader []interface{}, rows [][]interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Table", "c": []interface{}{caption, colAlign, relColWidth,
		colHeader, rows}}
}

func Div(attrs []interface{}, blocks []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Div", "c": []interface{}{attrs, blocks}}
}

func Null() map[string]interface{} {
	return map[string]interface{}{"t": "Null", "c": []interface{}{}}
}

// Constructors for inline elements

func Str(text string) map[string]interface{} {
	return map[string]interface{}{"t": "Str", "c": text}
}

func Emph(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Emph", "c": inlines}
}

func Strong(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Strong", "c": inlines}
}

func Strikeout(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Strikeout", "c": inlines}
}

func Superscript(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Superscript", "c": inlines}
}

func Subscript(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Subscript", "c": inlines}
}

func SmallCaps(inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "SmallCaps", "c": inlines}
}

func Quoted(quoteType interface{}, inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Quoted", "c": []interface{}{quoteType, inlines}}
}

func Cite(citations []interface{}, inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Cite", "c": []interface{}{citations, inlines}}
}

func Code(attrs []interface{}, str string) map[string]interface{} {
	return map[string]interface{}{"t": "Code", "c": []interface{}{attrs, str}}
}

func Space() map[string]interface{} {
	return map[string]interface{}{"t": "Space", "c": []interface{}{}}
}

func SoftBreak() map[string]interface{} {
	return map[string]interface{}{"t": "SoftBreak", "c": []interface{}{}}
}

func LineBreak() map[string]interface{} {
	return map[string]interface{}{"t": "LineBreak", "c": []interface{}{}}
}

func Math(mathType interface{}, str string) map[string]interface{} {
	return map[string]interface{}{"t": "Math", "c": []interface{}{mathType, str}}
}

func RawInline(format string, str string) map[string]interface{} {
	return map[string]interface{}{"t": "RawInline", "c": []interface{}{format, str}}
}

func Link(attrs []interface{}, inlines []interface{}, target []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Link", "c": []interface{}{attrs, inlines, target}}
}

func Image(attrs []interface{}, inlines []interface{}, target []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Image", "c": []interface{}{attrs, inlines, target}}
}

func Note(blocks []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Note", "c": blocks}
}

func Span(attrs []interface{}, inlines []interface{}) map[string]interface{} {
	return map[string]interface{}{"t": "Span", "c": []interface{}{attrs, inlines}}
}
