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
	s.Router.NotFoundHandler = http.HandlerFunc(notFound)

	s.Router.Handle("/", http.HandlerFunc(s.home)).Methods(http.MethodGet)
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

func (srv Server) home(w http.ResponseWriter, r *http.Request) {
	sid, err := srv.session.GetUserSID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	w.Write([]byte(sid))
}

func login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
