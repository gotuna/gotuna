package app

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/renderer"
	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//go:embed public/*
var embededPublic embed.FS

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

func NewServer(logger *log.Logger, s *session.Session, userRepository UserRepository) *Server {
	srv := &Server{}
	srv.session = s
	srv.userRepository = userRepository

	srv.Router = mux.NewRouter()
	srv.Router.NotFoundHandler = srv.notFound()

	srv.Router.Handle("/", srv.home()).Methods(http.MethodGet)
	srv.Router.Handle("/login", srv.login()).Methods(http.MethodGet)
	srv.Router.Handle("/login", srv.loginSubmit()).Methods(http.MethodPost)
	srv.Router.Handle("/register", srv.login()).Methods(http.MethodGet, http.MethodPost)

	//bad := func() http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		var x map[string]int
	//		x["y"] = 1 // will produce nil map panic
	//	})
	//}
	//srv.Router.Handle("/bad", bad())

	srv.Router.Use(middleware.Logger(logger))
	srv.Router.Use(middleware.AuthRedirector(srv.session))

	// serve files from the public directory
	srv.Router.PathPrefix("/public/").Handler(ServeFiles(embededPublic))

	return srv
}

func ServeFiles(filesystem fs.FS) http.Handler {
	fs := http.FS(filesystem)
	filesrv := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			//NotFoundHandler(w, r)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		//stat, _ := f.Stat()
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		filesrv.ServeHTTP(w, r)
	})
}

func (srv Server) home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := renderer.NewHTMLRenderer("home.html")
		t.Render(w, http.StatusOK)
	})
}

func (srv Server) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := renderer.NewHTMLRenderer("login.html")
		t.Render(w, http.StatusOK)
	})
}

func (srv Server) loginSubmit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		user, err := srv.userRepository.GetUserByEmail(email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
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
