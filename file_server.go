package gotuna

import (
	"net/http"
	"os"
	"path"
)

// ServeFiles returns a Handler that serves a HTTP requests with the file contents.
//
// You can pass 404 Handler to be served when the file is not found.
//
// It will not list directories, instead, the 404 Handler will be used.
func (app App) ServeFiles(notFound http.Handler) http.Handler {
	fs := http.FS(app.Static)
	fileapp := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			notFound.ServeHTTP(w, r)
			return
		}
		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			notFound.ServeHTTP(w, r)
			return
		}

		// TODO: ModTime doesn't work for embed?
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		fileapp.ServeHTTP(w, r)
	})
}
