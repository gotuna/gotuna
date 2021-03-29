package templating

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/util"
	"github.com/alcalbg/gotdd/views"
)

type TemplatingEngine interface {
	Render(w http.ResponseWriter, r *http.Request, patterns ...string)
	Set(key string, data interface{}) TemplatingEngine
	SetError(errorKey, description string) TemplatingEngine
	GetErrors() map[string]string
	MountFS(fs fs.FS) TemplatingEngine
}

func GetEngine(t i18n.Locale, ses *session.Session) TemplatingEngine {

	var funcs = template.FuncMap{
		"t": t.T,
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
		"static": func(file string) string {
			hash := "b1a2"
			return fmt.Sprintf("%s%s?%s", util.StaticPath, file, hash)
		},
	}

	return &nativeHtmlTemplates{
		fs:      views.EmbededViews,
		funcs:   funcs,
		Data:    make(map[string]interface{}),
		Errors:  make(map[string]string),
		session: ses,
	}
}

type nativeHtmlTemplates struct {
	fs      fs.FS
	funcs   template.FuncMap
	Data    map[string]interface{}
	Errors  map[string]string
	Request *http.Request
	session *session.Session
	Flashes []session.FlashMessage
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

	w.Header().Set("Content-type", util.ContentTypeHTML)

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
