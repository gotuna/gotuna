package app

import (
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/util"
)

type App struct {
	util.Options
}

func NewApp(options util.Options) http.Handler {

	app := &App{util.OptionsWithDefaults(options)}

	app.Router.NotFoundHandler = http.HandlerFunc(app.notFound)

	// middlewares for all routes
	app.Router.Use(middleware.Recoverer(app.Options))
	app.Router.Use(middleware.Logger(app.Options))
	// TODO: csrf middleware
	app.Router.Methods(http.MethodOptions)
	app.Router.Use(middleware.Cors())

	// logged in user
	user := app.Router.NewRoute().Subrouter()
	user.Use(middleware.Authenticate(app.Options, "/login"))
	user.Handle("/", app.home()).Methods(http.MethodGet)
	user.Handle("/profile", app.profile()).Methods(http.MethodGet, http.MethodPost)
	user.Handle("/logout", app.logout()).Methods(http.MethodPost)

	// guests
	auth := app.Router.NewRoute().Subrouter()
	auth.Use(middleware.RedirectIfAuthenticated(app.Options, "/"))
	auth.Handle("/login", app.login()).Methods(http.MethodGet, http.MethodPost)
	auth.Handle("/register", app.login()).Methods(http.MethodGet, http.MethodPost)

	// serve files from the static directory
	app.Router.PathPrefix(app.StaticPrefix).
		Handler(http.StripPrefix(app.StaticPrefix, app.serveFiles())).
		Methods(http.MethodGet)

	return app.Router
}

func (app App) serveFiles() http.Handler {
	fs := http.FS(app.FS)
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

		templating.GetEngine(app.Options).
			Set("message", app.Locale.T("Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func (app App) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := templating.GetEngine(app.Options)

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" {
			tmpl.SetError("email", app.Locale.T("This field is required"))
		}
		if password == "" {
			tmpl.SetError("password", app.Locale.T("This field is required"))
		}
		if len(tmpl.GetErrors()) > 0 {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		app.UserRepository.Set("email", email)
		app.UserRepository.Set("password", password)
		user, err := app.UserRepository.Authenticate()
		if err != nil {
			tmpl.SetError("email", app.Locale.T("Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		if err := app.Session.SetUserSID(w, r, user.SID); err != nil {
			app.errorHandler(err).ServeHTTP(w, r)
			return
		}

		if err := app.Session.Flash(w, r, session.NewFlash(app.Locale.T("Welcome"))); err != nil {
			app.errorHandler(err).ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (app App) logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Session.DestroySession(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func (app App) profile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templating.GetEngine(app.Options).
			Render(w, r, "app.html", "profile.html")
	})
}

func (app App) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	templating.GetEngine(app.Options).
		Set("title", app.Locale.T("Not found")).
		Render(w, r, "app.html", "4xx.html")
}

func (app App) errorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		templating.GetEngine(app.Options).
			Set("error", err).
			Set("stacktrace", string(debug.Stack())).
			Render(w, r, "app.html", "error.html")
	})
}
