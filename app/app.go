package app

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/static"
	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
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
	lang           i18n.Translator
}

func NewServer(logger *log.Logger, s *session.Session, userRepository UserRepository) *Server {

	srv := &Server{}
	srv.session = s
	srv.userRepository = userRepository
	srv.lang = i18n.NewTranslator(i18n.En) // TODO: move this to session/user/store

	srv.Router = mux.NewRouter()
	srv.Router.NotFoundHandler = http.HandlerFunc(srv.notFound)

	srv.Router.Handle("/", srv.home()).Methods(http.MethodGet)
	srv.Router.Handle("/login", srv.login()).Methods(http.MethodGet, http.MethodPost)
	srv.Router.Handle("/logout", srv.logout()).Methods(http.MethodPost)
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

	// serve files from the static directory
	srv.Router.PathPrefix(util.StaticPath).
		Handler(http.StripPrefix(util.StaticPath,
			srv.ServeFiles(static.EmbededStatic)))

	return srv
}

func (srv Server) ServeFiles(filesystem fs.FS) http.Handler {
	fs := http.FS(filesystem)
	filesrv := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			srv.notFound(w, r)
			return
		}
		stat, _ := f.Stat()
		if stat.IsDir() {
			srv.notFound(w, r) // do not show directory listing
			return
		}
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		filesrv.ServeHTTP(w, r)
	})
}

func (srv Server) home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		templating.GetEngine(srv.lang).
			Set("message", srv.lang.T("Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func (srv Server) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := templating.GetEngine(srv.lang)

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" {
			tmpl.AddError("email", srv.lang.T("This field is required"))
		}
		if password == "" {
			tmpl.AddError("password", srv.lang.T("This field is required"))
		}
		if len(tmpl.GetErrors()) > 0 {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		user, err := srv.userRepository.GetUserByEmail(email)
		if err != nil {
			tmpl.AddError("email", srv.lang.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			tmpl.AddError("email", srv.lang.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		srv.session.SetUserSID(w, r, user.SID)

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (srv Server) logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.session.DestroySession(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func (srv Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
