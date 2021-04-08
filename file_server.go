package gotdd

import (
	"net/http"
	"os"
	"path"
)

func (app App) ServeFiles() http.Handler {
	fs := http.FS(app.Static)
	fileapp := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		stat, _ := f.Stat()
		if stat.IsDir() {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		// TODO: ModTime doesn't work for embed?
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		fileapp.ServeHTTP(w, r)
	})
}
