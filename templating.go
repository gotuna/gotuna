package gotuna

import (
	"fmt"
	"html/template"
	"net/http"
)

const ContentTypeHTML = "text/html; charset=utf-8"

type ViewHelperFunc func(w http.ResponseWriter, r *http.Request) (string, interface{})

type TemplatingEngine interface {
	Render(w http.ResponseWriter, r *http.Request, patterns ...string)
	Set(key string, data interface{}) TemplatingEngine
	SetError(errorKey, description string) TemplatingEngine
	GetErrors() map[string]string
}

func (app App) NewTemplatingEngine() TemplatingEngine {
	return &nativeHtmlTemplates{
		app:    app,
		Data:   make(map[string]interface{}),
		Errors: make(map[string]string),
	}
}

type nativeHtmlTemplates struct {
	app     App
	Data    map[string]interface{}
	Errors  map[string]string
	Flashes []FlashMessage
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

	if t.app.Session != nil {
		t.Flashes = t.app.Session.Flashes(w, r)
	}

	tmpl := template.Must(
		template.New("app").
			Funcs(t.getHelpers(w, r)).
			ParseFS(t.app.ViewFiles, patterns...))

	w.Header().Set("Content-type", ContentTypeHTML)

	err := tmpl.Execute(w, t)
	if err != nil {
		panic(fmt.Errorf("error rendering %v", err))
	}
}

func (t nativeHtmlTemplates) getHelpers(w http.ResponseWriter, r *http.Request, patterns ...string) template.FuncMap {
	// default helpers
	fmap := template.FuncMap{
		"request": func() *http.Request {
			return r
		},
		"t": func(s string) string {
			locale := t.app.Session.GetUserLocale(r)
			return t.app.Locale.T(locale, s)
		},
		"static": func(file string) string {
			hash := "123" // TODO:
			return fmt.Sprintf("%s%s?%s", t.app.StaticPrefix, file, hash)
		},
		"currentUser": func() (User, error) {
			user, err := GetUserFromContext(r.Context())
			return user, err
		},
		"currentLocale": func() string {
			return t.app.Session.GetUserLocale(r)
		},
		"isGuest": func() bool {
			return t.app.Session.IsGuest(r)
		},
	}
	// add custom, user-defined helpers
	for _, v := range t.app.ViewHelpers {
		n, f := v(w, r)
		fmap[n] = f
	}
	return fmap
}
