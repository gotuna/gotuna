package renderer

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

const BaseDir = "views"
const LayoutFile = "app.html"
const ContentTypeHTML = "text/html; charset=utf-8"

//go:embed views/*
var embededViews embed.FS

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

func NewHTMLRenderer(filename string) Renderer {
	return &htmlRenderer{
		filename: filename,
		fs:       embededViews,
		Data:     make(map[string]interface{}),
	}
}

type Renderer interface {
	Render(w http.ResponseWriter, statusCode int)
	Set(key string, data interface{})
	Mount(fs fs.FS)
}

type htmlRenderer struct {
	filename string
	fs       fs.FS
	Data     map[string]interface{}
}

func (t *htmlRenderer) Set(key string, data interface{}) {
	t.Data[key] = data
}

func (t htmlRenderer) Render(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-type", ContentTypeHTML)
	w.WriteHeader(statusCode)

	tmpl := template.Must(template.New("app").Funcs(funcs).ParseFS(t.fs,
		fmt.Sprintf("%s/%s", BaseDir, LayoutFile),
		fmt.Sprintf("%s/%s", BaseDir, t.filename),
	),
	)

	err := tmpl.Execute(w, t)
	if err != nil {
		fmt.Println(err)
		panic("TODO")
	}
}

func (t *htmlRenderer) Mount(fs fs.FS) {
	t.fs = fs
}
