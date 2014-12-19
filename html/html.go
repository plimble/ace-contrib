package html

import (
	"github.com/flosch/pongo2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Context map[string]interface{}

type Options struct {
	Directory     string
	Extensions    []string
	IsDevelopment bool
}

type HTML struct {
	options *Options
	Set     *pongo2.TemplateSet
}

func New(options *Options) *HTML {
	html := &HTML{}
	if options.IsDevelopment {
		pongo2.DefaultSet.Debug = true
	}

	html.options = options
	html.Set = pongo2.DefaultSet
	html.Set.SetBaseDirectory(options.Directory)
	html.compileTemplates()
	return html
}

func (html *HTML) ExecW(name string, context map[string]interface{}, w http.ResponseWriter) {
	tpl, err := html.Set.FromCache(name)
	if err != nil {
		panic(err)
	}
	if err := tpl.ExecuteWriter(context, w); err != nil {
		panic(err)
	}
}

func (html *HTML) Exec(name string, context map[string]interface{}) string {
	tpl, err := html.Set.FromCache(name)
	if err != nil {
		panic(err)
	}
	result, err := tpl.Execute(context)
	if err != nil {
		panic(err)
	}
	return result
}

func (html *HTML) ExecByte(name string, context map[string]interface{}) []byte {
	tpl, err := html.Set.FromCache(name)
	if err != nil {
		panic(err)
	}

	result, err := tpl.ExecuteBytes(context)
	if err != nil {
		panic(err)
	}

	return result
}

func (html *HTML) compileTemplates() {
	dir := html.options.Directory

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := ""
		if strings.Index(rel, ".") != -1 {
			ext = "." + strings.Join(strings.Split(rel, ".")[1:], ".")
		}

		for _, extension := range html.options.Extensions {
			if ext == extension {
				name := (rel[0 : len(rel)-len(ext)])
				pongo2.Must(html.Set.FromCache(filepath.ToSlash(name) + ext))
				break
			}
		}

		return nil
	})
}
