package app

import (
	"log"
	"net/http"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Server struct {
	Router  *mux.Router
	session *session.Session
}

func NewServer(logger *log.Logger, sessionStore sessions.Store) *Server {
	s := &Server{}
	s.session = session.NewSession(sessionStore)

	s.Router = mux.NewRouter()
	s.Router.NotFoundHandler = s.notFound()

	s.Router.Handle("/", s.home()).Methods(http.MethodGet)
	s.Router.Handle("/login", s.login()).Methods(http.MethodGet)

	//bad := func() http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		var x map[string]int
	//		x["y"] = 1 // will produce nil map panic
	//	})
	//}
	//s.Router.Handle("/bad", bad())

	s.Router.Use(middleware.Logger(logger))
	s.Router.Use(middleware.AuthRedirector(s.session))

	return s
}

func (srv Server) home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (srv Server) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (srv Server) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
