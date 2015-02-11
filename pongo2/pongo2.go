package pongo2

import (
	"github.com/plimble/ace"
	"github.com/plimble/copter"
	"net/http"
)

type TemplateOptions struct {
	Directory     string
	Extensions    []string
	IsDevelopment bool
}

type pongo2 struct {
	copter *copter.Copter
}

func (p *pongo2) Render(w http.ResponseWriter, name string, data interface{}) {
	p.copter.ExecW(name, data.(map[string]interface{}), w)
}

func Pongo2(options *TemplateOptions) ace.Renderer {
	return &pongo2{
		copter: copter.New(&copter.Options{
			Directory:     options.Directory,
			Extensions:    options.Extensions,
			IsDevelopment: options.IsDevelopment,
		}),
	}
}
