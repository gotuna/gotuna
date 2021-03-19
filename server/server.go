package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

const get = http.MethodGet
const post = http.MethodPost
const DefaultError = "Whoops, something went wrong"

func NewServer() *Server {
	s := &Server{}
	s.Router = mux.NewRouter()
	s.Router.NotFoundHandler = http.HandlerFunc(notFound)

	s.Router.Handle("/", home()).Methods(get)
	s.Router.Handle("/login", login()).Methods(get)

	s.Router.Use(recoverer)

	return s
}

func home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func recoverer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(DefaultError))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
