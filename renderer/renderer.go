package renderer

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

const BaseDir = "views"
const LayoutFile = "app.html"
const ContentTypeHTML = "text/html; charset=utf-8"

//go:embed views/*
var embededViews embed.FS

//go:embed public/*
var embededPublic embed.FS

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

type Renderer interface {
	Render(w http.ResponseWriter, statusCode int)
	Set(key string, data interface{})
}

func NewHTMLRenderer(filename string, fs fs.FS) Renderer {
	if fs == nil {
		fs = embededViews
	}
	return &htmlRenderer{
		filename: filename,
		fs:       fs,
		Data:     make(map[string]interface{}),
	}
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

func ServeFiles(filesystem fs.FS) http.Handler {
	if filesystem == nil {
		filesystem = embededPublic
	}
	fs := http.FS(filesystem)
	filesrv := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			//NotFoundHandler(w, r)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		//stat, _ := f.Stat()
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		filesrv.ServeHTTP(w, r)
	})
}
