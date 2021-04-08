package fullapp

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gotuna/gotuna"
)

var User1 = gotuna.InMemoryUser{
	UniqueID: "123",
	Email:    "john@example.com",
	Name:     "John",
	Password: "pass123",
}

var User2 = gotuna.InMemoryUser{
	UniqueID: "456",
	Email:    "bob@example.com",
	Name:     "Bob",
	Password: "bobby5",
}

func NewUserRepository() gotuna.UserRepository {
	return gotuna.NewInMemoryUserRepository([]gotuna.InMemoryUser{
		User1,
		User2,
	})
}

func MakeApp(app gotuna.App) gotuna.App {

	if app.Logger == nil {
		app.Logger = log.New(io.Discard, "", 0)
	}

	if app.Locale == nil {
		app.Locale = gotuna.NewLocale(map[string]map[string]string{})
	}

	// custom view helpers
	app.ViewHelpers = []gotuna.ViewHelperFunc{
		func(w http.ResponseWriter, r *http.Request) (string, interface{}) {
			return "uppercase", func(s string) string {
				return strings.ToUpper(s)
			}
		},
		func(w http.ResponseWriter, r *http.Request) (string, interface{}) {
			return "csrf", func() template.HTML {
				return csrf.TemplateField(r)
			}
		},
	}

	app.Router = mux.NewRouter()

	// middlewares for all routes
	app.Router.Handle("/error", handlerError(app)).Methods(http.MethodGet, http.MethodPost)
	app.Router.Use(app.Recoverer("/error"))
	app.Router.Use(app.Logging())
	app.Router.Methods(http.MethodOptions)
	app.Router.Use(app.Cors())

	// for logged in users
	user := app.Router.NewRoute().Subrouter()
	user.Use(app.Authenticate("/login"))
	user.Use(app.StoreUserToContext())
	user.Handle("/", handlerHome(app)).Methods(http.MethodGet)
	user.Handle("/profile", handlerProfile(app)).Methods(http.MethodGet, http.MethodPost)
	user.Handle("/logout", handlerLogout(app)).Methods(http.MethodPost)
	user.Handle("/setlocale/{locale}", handlerChangeLocale(app)).Methods(http.MethodGet, http.MethodPost)

	// for guests
	auth := app.Router.NewRoute().Subrouter()
	auth.Use(app.RedirectIfAuthenticated("/"))
	auth.Handle("/login", handlerLogin(app)).Methods(http.MethodGet, http.MethodPost)

	// path prefix for static files
	// e.g. "/public" or "http://cdn.example.com/assets"
	app.StaticPrefix = strings.TrimRight(app.StaticPrefix, "/")

	// serve files from the static directory
	app.Router.PathPrefix(app.StaticPrefix).
		Handler(http.StripPrefix(app.StaticPrefix, app.ServeFiles(handlerNotFound(app)))).
		Methods(http.MethodGet)

	return app
}

func handlerHome(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewTemplatingEngine().
			Set("message", app.Locale.T(app.Session.GetUserLocale(r), "Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func handlerLogin(app gotuna.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl := app.NewTemplatingEngine()

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		user, err := app.UserRepository.Authenticate(w, r)
		if err != nil {
			tmpl.SetError("email", app.Locale.T(app.Session.GetUserLocale(r), "Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		if err := app.Session.SetUserID(w, r, user.GetID()); err != nil {
			app.Logger.Printf("%s %s %s %v", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, err)
			handlerError(app).ServeHTTP(w, r)
			return
		}

		app.Session.SetUserLocale(w, r, "en-US")

		flash(app, w, r, t(app, r, "Welcome"))

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func handlerLogout(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Session.Destroy(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func handlerProfile(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewTemplatingEngine().
			Render(w, r, "app.html", "profile.html")
	})
}

func handlerChangeLocale(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		app.Session.SetUserLocale(w, r, vars["locale"])
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func handlerNotFound(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		app.NewTemplatingEngine().
			Set("title", app.Locale.T(app.Session.GetUserLocale(r), "Not found")).
			Render(w, r, "404.html")
	})
}

func handlerError(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		app.NewTemplatingEngine().
			Set("error", "TODO"). // TODO: show error
			Set("stacktrace", string(debug.Stack())).
			Render(w, r, "error.html")
	})
}

func flash(app gotuna.App, w http.ResponseWriter, r *http.Request, msg string) {
	if err := app.Session.Flash(w, r, gotuna.NewFlash(msg)); err != nil {
		app.Logger.Printf("%s %s %s %v", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, err)
	}
}

func t(app gotuna.App, r *http.Request, s string) string {
	return app.Locale.T(app.Session.GetUserLocale(r), s)
}
