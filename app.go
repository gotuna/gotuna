package gotdd

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
)

type App struct {
	Logger         *log.Logger
	Router         *mux.Router
	FS             fs.FS
	Session        *Session
	UserRepository UserRepository
	StaticPrefix   string
	Locale         Locale
}

func NewApp(app App) App {

	if app.Logger == nil {
		app.Logger = log.New(os.Stdout, "", 0)
	}

	if app.Locale == nil {
		app.Locale = NewLocale(Translations)
	}

	app.Router = mux.NewRouter()
	app.Router.NotFoundHandler = http.HandlerFunc(app.notFound)

	// middlewares for all routes
	app.Router.Use(app.Recoverer())
	app.Router.Use(app.Logging())
	// TODO: csrf middleware
	app.Router.Methods(http.MethodOptions)
	app.Router.Use(app.Cors())

	// logged in user
	user := app.Router.NewRoute().Subrouter()
	user.Use(app.Authenticate("/login"))
	user.Handle("/", app.home()).Methods(http.MethodGet)
	user.Handle("/profile", app.profile()).Methods(http.MethodGet, http.MethodPost)
	user.Handle("/logout", app.logout()).Methods(http.MethodPost)

	// guests
	auth := app.Router.NewRoute().Subrouter()
	auth.Use(app.RedirectIfAuthenticated("/"))
	auth.Handle("/login", app.login()).Methods(http.MethodGet, http.MethodPost)
	auth.Handle("/register", app.login()).Methods(http.MethodGet, http.MethodPost)

	// path prefix for static files
	// e.g. "/public" or "http://cdn.example.com/assets"
	app.StaticPrefix = strings.TrimRight(app.StaticPrefix, "/")

	// serve files from the static directory
	app.Router.PathPrefix(app.StaticPrefix).
		Handler(http.StripPrefix(app.StaticPrefix, app.serveFiles())).
		Methods(http.MethodGet)

	return app
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

		app.GetEngine().
			Set("message", app.Locale.T("en-US", "Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func (app App) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := app.GetEngine()

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" {
			tmpl.SetError("email", app.Locale.T("en-US", "This field is required"))
		}
		if password == "" {
			tmpl.SetError("password", app.Locale.T("en-US", "This field is required"))
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
			tmpl.SetError("email", app.Locale.T("en-US", "Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		if err := app.Session.SetUserSID(w, r, user.SID); err != nil {
			app.errorHandler(err).ServeHTTP(w, r)
			return
		}

		if err := app.Session.Flash(w, r, NewFlash(app.Locale.T("en-US", "Welcome"))); err != nil {
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
		app.GetEngine().
			Render(w, r, "app.html", "profile.html")
	})
}

func (app App) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	app.GetEngine().
		Set("title", app.Locale.T("en-US", "Not found")).
		Render(w, r, "app.html", "4xx.html")
}

func (app App) errorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		app.GetEngine().
			Set("error", err).
			Set("stacktrace", string(debug.Stack())).
			Render(w, r, "app.html", "error.html")
	})
}
