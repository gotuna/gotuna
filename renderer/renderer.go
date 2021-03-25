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

func NewHTMLRenderer(translator lang.Translator, patterns ...string) Renderer {
	return &htmlRenderer{
		patterns: patterns,
		fs:       views.EmbededViews,
		Data:     make(map[string]interface{}),
		Lang:     translator,
	}
}

type Renderer interface {
	Render(w http.ResponseWriter, statusCode int) error
	Set(key string, data interface{})
	Mount(fs fs.FS)
}

type htmlRenderer struct {
	patterns []string
	fs       fs.FS
	Data     map[string]interface{}
	Lang     lang.Translator
}

func (t *htmlRenderer) Set(key string, data interface{}) {
	t.Data[key] = data
}

func (t htmlRenderer) Render(w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-type", ContentTypeHTML)
	w.WriteHeader(statusCode)

	tmpl := template.Must(
		template.New("app").
			Funcs(funcs).
			ParseFS(t.fs, t.patterns...))

	err := tmpl.Execute(w, t)
	if err != nil {
		return fmt.Errorf("error rendering %v", err)
	}

	return nil
}

func (t *htmlRenderer) Mount(fs fs.FS) {
	t.fs = fs
}
