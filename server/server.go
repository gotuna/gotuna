package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct{}

func NewServer() http.Handler {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	router.Handle("/", http.HandlerFunc(homeHandler))
	router.Handle("/login", http.HandlerFunc(loginHandler))
	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
