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
	Render(w http.ResponseWriter, patterns ...string) error
	Set(key string, data interface{}) TemplatingEngine
	Mount(fs fs.FS) TemplatingEngine
}

func GetNativeTemplatingEngine(translator i18n.Translator) TemplatingEngine {

	var funcs = template.FuncMap{
		"lang": translator.Translate,
		"uppercase": func(v string) string {
			return strings.ToUpper(v)
		},
	}

	return &nativeHtmlTemplates{
		fs:    views.EmbededViews,
		funcs: funcs,
		Data:  make(map[string]interface{}),
	}
}

type nativeHtmlTemplates struct {
	fs    fs.FS
	funcs template.FuncMap
	Data  map[string]interface{}
}

func (t *nativeHtmlTemplates) Set(key string, data interface{}) TemplatingEngine {
	t.Data[key] = data
	return t
}

func (t nativeHtmlTemplates) Render(w http.ResponseWriter, patterns ...string) error {
	w.Header().Set("Content-type", util.ContentTypeHTML)

	tmpl := template.Must(
		template.New("app").
			Funcs(t.funcs).
			ParseFS(t.fs, patterns...))

	err := tmpl.Execute(w, t)
	if err != nil {
		return fmt.Errorf("error rendering %v", err)
	}

	return nil
}

func (t *nativeHtmlTemplates) Mount(fs fs.FS) TemplatingEngine {
	t.fs = fs
	return t
}
