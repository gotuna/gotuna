package templating

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/util"
	"github.com/alcalbg/gotdd/views"
)

type TemplatingEngine interface {
	Render(w http.ResponseWriter, r *http.Request, patterns ...string)
	Set(key string, data interface{}) TemplatingEngine
	SetError(errorKey, description string) TemplatingEngine
	GetErrors() map[string]string
	Mount(fs fs.FS) TemplatingEngine
}

func GetEngine(lang i18n.Translator) TemplatingEngine {

	var funcs = template.FuncMap{
		"lang": lang.T,
		"uppercase": func(v string) string {
			return strings.ToUpper(v)
		},
	}

	return &nativeHtmlTemplates{
		fs:     views.EmbededViews,
		funcs:  funcs,
		Data:   make(map[string]interface{}),
		Errors: make(map[string]string),
	}
}

type nativeHtmlTemplates struct {
	fs      fs.FS
	funcs   template.FuncMap
	Data    map[string]interface{}
	Errors  map[string]string
	Request *http.Request
	Ver     string
}

func (t *nativeHtmlTemplates) Set(key string, data interface{}) TemplatingEngine {
	t.Data[key] = data
	return t
}

func (t *nativeHtmlTemplates) SetError(errorKey, description string) TemplatingEngine {
	t.Errors[errorKey] = description
	return t
}

func (t nativeHtmlTemplates) GetErrors() map[string]string {
	return t.Errors
}

func (t *nativeHtmlTemplates) Render(w http.ResponseWriter, r *http.Request, patterns ...string) {

	w.Header().Set("Content-type", util.ContentTypeHTML)

	t.Request = r
	t.Ver = "22" // TODO: fix this

	tmpl := template.Must(
		template.New("app").
			Funcs(t.funcs).
			ParseFS(t.fs, patterns...))

	err := tmpl.Execute(w, t)
	if err != nil {
		panic(fmt.Errorf("error rendering %v", err))
	}
}

func (t *nativeHtmlTemplates) Mount(fs fs.FS) TemplatingEngine {
	t.fs = fs
	return t
}
