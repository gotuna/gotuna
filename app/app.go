package app

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/models"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	Mux            *mux.Router
	session        *session.Session
	userRepository models.UserRepository
	lang           i18n.Translator
	fs             fs.FS
}

func NewServer(logger *log.Logger, fs fs.FS, s *session.Session, userRepository models.UserRepository) *Server {

	srv := &Server{}
	srv.session = s
	srv.fs = fs
	srv.userRepository = userRepository
	srv.lang = i18n.NewTranslator(i18n.En) // TODO: move this to session/user/store

	srv.Mux = mux.NewRouter()
	srv.Mux.NotFoundHandler = http.HandlerFunc(srv.notFound)

	srv.Mux.Handle("/", srv.home()).Methods(http.MethodGet)
	srv.Mux.Handle("/login", srv.login()).Methods(http.MethodGet, http.MethodPost)
	srv.Mux.Handle("/logout", srv.logout()).Methods(http.MethodPost)
	srv.Mux.Handle("/profile", srv.profile()).Methods(http.MethodGet, http.MethodPost)
	srv.Mux.Handle("/register", srv.login()).Methods(http.MethodGet, http.MethodPost)

	//bad := func() http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		var x map[string]int
	//		x["y"] = 1 // will produce nil map panic
	//	})
	//}
	//srv.Mux.Handle("/bad", bad())

	srv.Mux.Use(middleware.Logger(logger))
	srv.Mux.Use(middleware.AuthRedirector(srv.session))

	// serve files from the static directory
	srv.Mux.PathPrefix(util.StaticPath).
		Handler(http.StripPrefix(util.StaticPath,
			srv.serveFiles()))

	return srv
}

func (srv Server) serveFiles() http.Handler {
	fs := http.FS(srv.fs)
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

		templating.GetEngine(srv.lang, srv.session).
			Set("message", srv.lang.T("Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func (srv Server) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := templating.GetEngine(srv.lang, srv.session)

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" {
			tmpl.SetError("email", srv.lang.T("This field is required"))
		}
		if password == "" {
			tmpl.SetError("password", srv.lang.T("This field is required"))
		}
		if len(tmpl.GetErrors()) > 0 {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		user, err := srv.userRepository.GetUserByEmail(email)
		if err != nil {
			tmpl.SetError("email", srv.lang.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			tmpl.SetError("email", srv.lang.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		if err := srv.session.SetUserSID(w, r, user.SID); err != nil {
			srv.errorHandler(err).ServeHTTP(w, r)
			return
		}

		if err := srv.session.AddFlash(w, r, srv.lang.T("Welcome"), "is-success", true); err != nil {
			srv.errorHandler(err).ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (srv Server) logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.session.DestroySession(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func (srv Server) profile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templating.GetEngine(srv.lang, srv.session).
			Render(w, r, "app.html", "profile.html")
	})
}

func (srv Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	templating.GetEngine(srv.lang, srv.session).
		Set("title", srv.lang.T("Not found")).
		Render(w, r, "app.html", "4xx.html")
}

func (srv Server) errorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		templating.GetEngine(i18n.NewTranslator(i18n.En), nil). // TODO lang
									Set("error", err).
									Set("stacktrace", string(debug.Stack())).
									Render(w, r, "app.html", "error.html")
	})
}
