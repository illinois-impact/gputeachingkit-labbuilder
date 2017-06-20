package pandoc

import (
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	pf "github.com/webgpu/gputeachingkit-labbuilder/pkg/pandocfilter"
	"golang.org/x/net/context"
)

type buildOptions struct {
	From       string
	Template   string
	Variables  map[string]string
	Type       string
	Standalone bool
}

var BuildOptions = buildOptions{}

func BuildOptionsFilter(k string, v interface{}, format string, meta interface{}) interface{} {
	if _, ok := ctx.Value("VisitedBuildOptionsFilter").(bool); ok {
		return nil
	}
	ctx = context.WithValue(ctx, "VisitedBuildOptionsFilter", true)

	info, ok := meta.(map[string]interface{})
	if !ok {
		pp.Println(meta)
	}

	if _, ok := info["build"]; !ok {
		logrus.Error("Cannot find document build in title block.\n")
		return nil
	}
	metamap := info["build"].(map[string]interface{})
	if _, ok := metamap["c"]; !ok {
		logrus.Error("Invalid document build format in title block.\n")
		return nil
	}
	build := metamap["c"].(map[string]interface{})

	if from, ok := build["from"]; ok {
		BuildOptions.From = pf.Stringify(from)
	}
	if ty, ok := build["type"]; ok {
		BuildOptions.Type = pf.Stringify(ty)
	}
	if template, ok := build["template"]; ok {
		BuildOptions.Template = pf.Stringify(template)
	}
	if standalone, ok := build["standalone"]; ok {
		BuildOptions.Standalone = standalone.(map[string]interface{})["c"].(bool)
	}
	if variables, ok := build["variables"]; ok {
		BuildOptions.Variables = make(map[string]string)
		for k, v := range variables.(map[string]interface{})["c"].(map[string]interface{}) {
			BuildOptions.Variables[k] = pf.Stringify(v)
		}
	}

	return nil

}

func init() {
	AddFilter(BuildOptionsFilter)
}
