package render

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/lang"
	"github.com/gorilla/csrf"
)

const defaultContentType = "text/html; charset=utf-8"

//go:embed views/*
var files embed.FS

//go:embed public/*
var public embed.FS

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

type Template struct {
	Filename    string
	Data        map[string]interface{}
	Errors      map[string]string
	Request     *http.Request
	StatusCode  int
	ContentType string
	Lang        lang.Translator
	Ver         int
}

func NewTemplate(filename string) *Template {

	lang.InitTranslator(lang.En)

	return &Template{
		Filename:    filename,
		ContentType: defaultContentType,
		StatusCode:  200,
		Data:        make(map[string]interface{}),
		Errors:      make(map[string]string),
		Lang:        lang.Lang,
		Ver:         56,
	}
}

// Set new view variable
func (t *Template) Set(key string, data interface{}) {
	t.Data[key] = data
}

func (t *Template) Render(w http.ResponseWriter, r *http.Request, code int) {

	// not instantiated with new?
	if t.Data == nil || t.Errors == nil {
		return
	}

	t.Request = r
	t.StatusCode = code

	t.Data["csrf"] = csrf.TemplateField(r)

	w.WriteHeader(t.StatusCode)
	w.Header().Set("Content-Type", t.ContentType)

	tmpl := template.Must(template.New("app").Funcs(funcs).ParseFS(files,
		"views/app.html",
		"views/"+t.Filename),
	)

	err := tmpl.Execute(w, t)
	if err != nil {
		fmt.Println(err)
		panic("TODO")
	}
}
