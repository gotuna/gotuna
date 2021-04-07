package gotdd

import (
	"fmt"
	"html/template"
	"net/http"
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

	if t.app.Session != nil {
		t.Flashes, _ = t.app.Session.Flashes(w, r)
		t.IsGuest = t.app.Session.IsGuest(r)
		t.UserLocale = t.app.Session.GetUserLocale(r)
	}

	t.Request = r

	tmpl := template.Must(
		template.New("app").
			Funcs(t.getTemplateFuncMap()).
			ParseFS(t.app.Views, patterns...))

	w.Header().Set("Content-type", ContentTypeHTML)

	err := tmpl.Execute(w, t)
	if err != nil {
		panic(fmt.Errorf("error rendering %v", err))
	}
}

func (t nativeHtmlTemplates) getTemplateFuncMap() template.FuncMap {
	// default helpers
	fmap := template.FuncMap{
		"t": func(s string) string {
			return t.app.Locale.T(t.UserLocale, s)
		},
		"static": func(file string) string {
			hash := "123" // TODO:
			return fmt.Sprintf("%s%s?%s", t.app.StaticPrefix, file, hash)
		},
		"currentuser": func() User {
			user, _ := GetUserFromContext(t.Request.Context())
			return user
		},
	}
	// add custom, user-defined helpers
	for k, v := range t.app.ViewsFuncMap {
		fmap[k] = v
	}
	return fmap
}
