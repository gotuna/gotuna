package app

import (
	"log"
	"net/http"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer(logger *log.Logger) *Server {
	s := &Server{}

	s.Router = mux.NewRouter()
	s.Router.NotFoundHandler = http.HandlerFunc(notFound)

	s.Router.Handle("/", home()).Methods(http.MethodGet)
	s.Router.Handle("/login", login()).Methods(http.MethodGet)

	//bad := func() http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		var x map[string]int
	//		x["y"] = 1 // will produce nil map panic
	//	})
	//}
	//s.Router.Handle("/bad", bad())

	s.Router.Use(middleware.Logger(logger))

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
