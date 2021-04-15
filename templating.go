package gotuna

import (
	"fmt"
	"html/template"
	"net/http"
)

// ContentTypeHTML is a standard Content-Type header for HTML response.
const ContentTypeHTML = "text/html; charset=utf-8"

// ViewHelperFunc is a view helper that can be used in any template file.
type ViewHelperFunc func(w http.ResponseWriter, r *http.Request) (string, interface{})

// TemplatingEngine used for rendering response.
type TemplatingEngine interface {
	Render(w http.ResponseWriter, r *http.Request, patterns ...string)
	Set(key string, data interface{}) TemplatingEngine
	SetError(errorKey, description string) TemplatingEngine
	GetErrors() map[string]string
}

// NewTemplatingEngine is a constructor that returns a native HTML templating engine.
func (app App) NewTemplatingEngine() TemplatingEngine {
	return &nativeHTML{
		app:    app,
		Data:   make(map[string]interface{}),
		Errors: make(map[string]string),
	}
}

type nativeHTML struct {
	app     App
	Data    map[string]interface{}
	Errors  map[string]string
	Flashes []FlashMessage
}

func (t *nativeHTML) Set(key string, data interface{}) TemplatingEngine {
	t.Data[key] = data
	return t
}

func (t *nativeHTML) SetError(errorKey, description string) TemplatingEngine {
	t.Errors[errorKey] = description
	return t
}

func (t nativeHTML) GetErrors() map[string]string {
	return t.Errors
}

func (t *nativeHTML) Render(w http.ResponseWriter, r *http.Request, patterns ...string) {

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

func (t nativeHTML) getHelpers(w http.ResponseWriter, r *http.Request, patterns ...string) template.FuncMap {
	// default helpers
	fmap := template.FuncMap{
		"request": func() *http.Request {
			return r
		},
		"t": func(s string) string {
			locale := t.app.Session.GetUserLocale(r)
			return t.app.Locale.T(locale, s)
		},
		"tp": func(s string, n int) string {
			locale := t.app.Session.GetUserLocale(r)
			return t.app.Locale.TP(locale, s, n)
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
