package gotdd

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/views"
)

const ContentTypeHTML = "text/html; charset=utf-8"

type TemplatingEngine interface {
	Render(w http.ResponseWriter, r *http.Request, patterns ...string)
	Set(key string, data interface{}) TemplatingEngine
	SetError(errorKey, description string) TemplatingEngine
	GetErrors() map[string]string
	MountFS(fs fs.FS) TemplatingEngine
}

func (app App) GetEngine() TemplatingEngine {

	translator := func(s string) string {
		return app.Locale.T("en-US", s) // TODO: set per user
	}

	var funcs = template.FuncMap{
		"t": translator,
		"static": func(file string) string {
			hash := "b1a2"
			return fmt.Sprintf("%s%s?%s", app.StaticPrefix, file, hash)
		},
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	return &nativeHtmlTemplates{
		fs:      views.EmbededViews,
		funcs:   funcs,
		Data:    make(map[string]interface{}),
		Errors:  make(map[string]string),
		session: app.Session,
	}
}

type nativeHtmlTemplates struct {
	fs      fs.FS
	funcs   template.FuncMap
	Data    map[string]interface{}
	Errors  map[string]string
	Request *http.Request
	session *Session
	Flashes []FlashMessage
	IsGuest bool
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

	w.Header().Set("Content-type", ContentTypeHTML)

	if t.session != nil {
		t.Flashes, _ = t.session.Flashes(w, r)
		t.IsGuest = t.session.IsGuest(r)
	}

	t.Request = r

	tmpl := template.Must(
		template.New("app").
			Funcs(t.funcs).
			ParseFS(t.fs, patterns...))

	err := tmpl.Execute(w, t)
	if err != nil {
		panic(fmt.Errorf("error rendering %v", err))
	}
}

func (t *nativeHtmlTemplates) MountFS(fs fs.FS) TemplatingEngine {
	t.fs = fs
	return t
}
