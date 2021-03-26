package renderer

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/lang"
	"github.com/alcalbg/gotdd/views"
)

const ContentTypeHTML = "text/html; charset=utf-8"

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

func NewHTMLRenderer(translator lang.Translator) Renderer {
	return &htmlRenderer{
		fs:   views.EmbededViews,
		Data: make(map[string]interface{}),
		Lang: translator,
	}
}

type Renderer interface {
	Render(w http.ResponseWriter, patterns ...string) error
	Set(key string, data interface{}) Renderer
	Mount(fs fs.FS) Renderer
}

type htmlRenderer struct {
	fs   fs.FS
	Data map[string]interface{}
	Lang lang.Translator
}

func (t *htmlRenderer) Set(key string, data interface{}) Renderer {
	t.Data[key] = data
	return t
}

func (t htmlRenderer) Render(w http.ResponseWriter, patterns ...string) error {
	w.Header().Set("Content-type", ContentTypeHTML)

	tmpl := template.Must(
		template.New("app").
			Funcs(funcs).
			ParseFS(t.fs, patterns...))

	err := tmpl.Execute(w, t)
	if err != nil {
		return fmt.Errorf("error rendering %v", err)
	}

	return nil
}

func (t *htmlRenderer) Mount(fs fs.FS) Renderer {
	t.fs = fs
	return t
}
