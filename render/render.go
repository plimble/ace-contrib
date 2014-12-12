package render

import (
	"gopkg.in/unrolled/render.v1"
	"html"
	"html/template"
	"reflect"
	"time"
)

type Render struct {
	*render.Render
}

type Options render.Options

func New(options Options) *Render {
	options.Funcs = append(options.Funcs, template.FuncMap{"esp": esp})
	options.Funcs = append(options.Funcs, template.FuncMap{"length": length})
	options.Funcs = append(options.Funcs, template.FuncMap{"date": date})
	options.Funcs = append(options.Funcs, template.FuncMap{"ISO8601": ISO8601})

	return &Render{render.New(render.Options(options))}
}

func esp(s string) string {
	return html.EscapeString(s)
}

func length(v interface{}) int {
	return reflect.TypeOf(v).Len()
}

func ISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}

func date(format string, t time.Time) string {
	return t.Format(format)
}
