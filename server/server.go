package server

import "net/http"

type Server struct{}

func (srv Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewServer() http.Handler {
	return Server{}
}
