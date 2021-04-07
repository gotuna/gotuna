package gotdd

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const ContentTypeHTML = "text/html; charset=utf-8"

type TemplatingEngine interface {
	Render(w http.ResponseWriter, r *http.Request, patterns ...string)
	Set(key string, data interface{}) TemplatingEngine
	SetError(errorKey, description string) TemplatingEngine
	GetErrors() map[string]string
}

func (app App) NewNativeTemplatingEngine() TemplatingEngine {
	return &nativeHtmlTemplates{
		app:    app,
		Data:   make(map[string]interface{}),
		Errors: make(map[string]string),
	}
}

type nativeHtmlTemplates struct {
	app        App
	Data       map[string]interface{}
	Errors     map[string]string
	Request    *http.Request
	Flashes    []FlashMessage
	UserLocale string
	IsGuest    bool
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

	if user, err := GetUser(r.Context()); err == nil {
		t.Set("currentuser", user)
	}

	if t.app.Session != nil {
		t.Flashes, _ = t.app.Session.Flashes(w, r)
		t.IsGuest = t.app.Session.IsGuest(r)
		t.UserLocale = t.app.Session.GetUserLocale(r)
	}

	t.Request = r

	tmpl := template.Must(
		template.New("app").
			Funcs(getTemplateFuncMap(t.app, t.UserLocale)).
			ParseFS(t.app.Views, patterns...))

	w.Header().Set("Content-type", ContentTypeHTML)

	err := tmpl.Execute(w, t)
	if err != nil {
		panic(fmt.Errorf("error rendering %v", err))
	}
}

func getTemplateFuncMap(app App, userLocale string) template.FuncMap {
	return template.FuncMap{
		"t": func(s string) string {
			return app.Locale.T(userLocale, s)
		},
		"static": func(file string) string {
			hash := "123" // TODO:
			return fmt.Sprintf("%s%s?%s", app.StaticPrefix, file, hash)
		},
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
	}
}
