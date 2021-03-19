package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct{}

func NewServer() http.Handler {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)

	router.Handle("/", http.HandlerFunc(home))
	router.Handle("/login", http.HandlerFunc(login))

	return router
}

func home(w http.ResponseWriter, r *http.Request) {
}

func login(w http.ResponseWriter, r *http.Request) {
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
