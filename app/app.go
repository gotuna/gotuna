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

type App struct {
	session        *session.Session
	userRepository models.UserRepository
	locale         i18n.Locale
	fs             fs.FS
}

func NewApp(logger *log.Logger, fs fs.FS, s *session.Session, userRepository models.UserRepository) http.Handler {

	app := &App{}
	app.session = s
	app.fs = fs
	app.userRepository = userRepository
	app.locale = i18n.NewLocale(i18n.En) // TODO: move this to session/user/store

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(app.notFound)

	router.Handle("/", app.home()).Methods(http.MethodGet)
	router.Handle("/login", app.login()).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/logout", app.logout()).Methods(http.MethodPost)
	router.Handle("/profile", app.profile()).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/register", app.login()).Methods(http.MethodGet, http.MethodPost)

	router.Use(middleware.Logger(logger))
	router.Use(middleware.AuthRedirector(app.session))

	// serve files from the static directory
	router.PathPrefix(util.StaticPath).
		Handler(http.StripPrefix(util.StaticPath,
			app.serveFiles()))

	return router
}

func (app App) serveFiles() http.Handler {
	fs := http.FS(app.fs)
	fileapp := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			app.notFound(w, r)
			return
		}
		stat, _ := f.Stat()
		if stat.IsDir() {
			app.notFound(w, r) // do not show directory listing
			return
		}

		// TODO: ModTime doesn't work for embed?
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		fileapp.ServeHTTP(w, r)
	})
}

func (app App) home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		templating.GetEngine(app.locale, app.session).
			Set("message", app.locale.T("Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func (app App) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := templating.GetEngine(app.locale, app.session)

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" {
			tmpl.SetError("email", app.locale.T("This field is required"))
		}
		if password == "" {
			tmpl.SetError("password", app.locale.T("This field is required"))
		}
		if len(tmpl.GetErrors()) > 0 {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		user, err := app.userRepository.GetUserByEmail(email)
		if err != nil {
			tmpl.SetError("email", app.locale.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			tmpl.SetError("email", app.locale.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		if err := app.session.SetUserSID(w, r, user.SID); err != nil {
			app.errorHandler(err).ServeHTTP(w, r)
			return
		}

		if err := app.session.AddFlash(w, r, app.locale.T("Welcome"), "is-success", true); err != nil {
			app.errorHandler(err).ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (app App) logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.session.DestroySession(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func (app App) profile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templating.GetEngine(app.locale, app.session).
			Render(w, r, "app.html", "profile.html")
	})
}

func (app App) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	templating.GetEngine(app.locale, app.session).
		Set("title", app.locale.T("Not found")).
		Render(w, r, "app.html", "4xx.html")
}

func (app App) errorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		templating.GetEngine(i18n.NewLocale(i18n.En), nil). // TODO lang
									Set("error", err).
									Set("stacktrace", string(debug.Stack())).
									Render(w, r, "app.html", "error.html")
	})
}
