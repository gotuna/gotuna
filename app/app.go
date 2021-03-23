package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	SID          string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	GetUserByEmail(email string) (User, error)
}

type Server struct {
	Router         *mux.Router
	session        *session.Session
	userRepository UserRepository
}

func NewServer(logger *log.Logger, gs *sessions.Session, userRepository UserRepository) *Server {
	s := &Server{}
	s.session = &session.Session{Store: gs.Store()}
	s.userRepository = userRepository

	s.Router = mux.NewRouter()
	s.Router.NotFoundHandler = s.notFound()

	s.Router.Handle("/", s.home()).Methods(http.MethodGet)
	s.Router.Handle("/login", s.login()).Methods(http.MethodGet, http.MethodPost)
	s.Router.Handle("/register", s.login()).Methods(http.MethodGet, http.MethodPost)

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
		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		user, err := srv.userRepository.GetUserByEmail(email)
		if err != nil {
			// TODO
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			// TODO - login failed
			return
		}

		// user is ok, save to session
		srv.session.SetUserSID(w, r, user.SID)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	})
}

func (srv Server) notFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
