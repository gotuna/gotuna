package fullapp

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gotuna/gotuna"
)

// App is the main dependency store.
type App struct {
	gotuna.App
	Csrf gotuna.MiddlewareFunc // app-specific config
}

// MakeApp creates and configures the App
func MakeApp(app App) App {

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

	// middlewares for all routes
	app.Router.Handle("/error", handlerError(app)).Methods(http.MethodGet, http.MethodPost)
	app.Router.Use(app.Recoverer("/error"))
	app.Router.Use(app.Logging())
	app.Router.Use(app.StoreParamsToContext())
	app.Router.Use(app.StoreUserToContext())
	app.Router.Use(app.Csrf)
	app.Router.Methods(http.MethodOptions)
	app.Router.Use(app.Cors())

	// for logged in users
	user := app.Router.NewRoute().Subrouter()
	user.Use(app.Authenticate("/login"))
	user.Handle("/", handlerHome(app)).Methods(http.MethodGet)
	user.Handle("/profile", handlerProfile(app)).Methods(http.MethodGet, http.MethodPost)
	user.Handle("/adduser", handlerAddUser(app)).Methods(http.MethodGet, http.MethodPost)
	user.Handle("/logout", handlerLogout(app)).Methods(http.MethodPost)
	user.Handle("/setlocale/{locale}", handlerChangeLocale(app)).Methods(http.MethodGet, http.MethodPost)

	// for guests
	guests := app.Router.NewRoute().Subrouter()
	guests.Use(app.RedirectIfAuthenticated("/"))
	guests.Handle("/login", handlerLogin(app)).Methods(http.MethodGet, http.MethodPost)

	// public json api
	api := app.Router.NewRoute().Subrouter()
	api.Handle("/api/getcars", handlerGetCars(app)).Methods(http.MethodGet)

	// path prefix for static files
	// e.g. "/public" or "http://cdn.example.com/assets"
	app.StaticPrefix = strings.TrimRight(app.StaticPrefix, "/")

	// serve files from the static directory
	app.Router.PathPrefix(app.StaticPrefix).
		Handler(http.StripPrefix(app.StaticPrefix, app.ServeFiles(handlerNotFound(app)))).
		Methods(http.MethodGet)

	return app
}

func handlerHome(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewTemplatingEngine().
			Set("message", app.Locale.T(app.Session.GetLocale(r), "Home")).
			Set("users", app.UserRepository.(*gotuna.InMemoryUserRepository).Users).
			Render(w, r, "app.html", "home.html")
	})
}

func handlerLogin(app App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl := app.NewTemplatingEngine()

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		user, err := app.UserRepository.Authenticate(w, r)
		if err != nil {
			tmpl.SetError("email", app.Locale.T(app.Session.GetLocale(r), "Login failed, please try again"))
			w.WriteHeader(http.StatusUnauthorized)
			tmpl.Render(w, r, "app.html", "login.html")
			return
		}

		// user is ok, save to session
		_ = app.Session.SetUserID(w, r, user.GetID())
		_ = app.Session.SetLocale(w, r, "en-US")

		// log this event
		app.Logger.Printf(
			"%s ### User logged in: %s",
			time.Now().Format(time.RFC3339),
			user.(gotuna.InMemoryUser).Name,
		)

		flash(app, w, r, t(app, r, "Welcome"))

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func handlerLogout(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, _ := gotuna.GetUserFromContext(r.Context())

		// log this event
		app.Logger.Printf(
			"%s ### User logged out: %s",
			time.Now().Format(time.RFC3339),
			user.(gotuna.InMemoryUser).Name,
		)

		_ = app.Session.Destroy(w, r)

		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func handlerProfile(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewTemplatingEngine().
			Render(w, r, "app.html", "profile.html")
	})
}

func handlerAddUser(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := app.NewTemplatingEngine()

		ctx := r.Context()
		repo := app.UserRepository.(*gotuna.InMemoryUserRepository)

		if r.Method == http.MethodGet {
			tmpl.Render(w, r, "app.html", "adduser.html")
			return
		}

		err := repo.AddUser(gotuna.InMemoryUser{
			ID:       gotuna.GetParam(ctx, "id"),
			Name:     gotuna.GetParam(ctx, "name"),
			Email:    gotuna.GetParam(ctx, "email"),
			Password: gotuna.GetParam(ctx, "password"),
		})

		if err != nil {
			_ = app.Session.Flash(w, r, gotuna.FlashMessage{
				Message:   t(app, r, "Error"),
				Kind:      "danger",
				AutoClose: true,
			})
			tmpl.Render(w, r, "app.html", "adduser.html")
			return
		}

		flash(app, w, r, t(app, r, "Success"))

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func handlerChangeLocale(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := gotuna.GetParam(r.Context(), "locale")
		_ = app.Session.SetLocale(w, r, locale)
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func handlerGetCars(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cars := []string{
			"Ford",
			"Škoda",
			"BMW",
		}

		json.NewEncoder(w).Encode(cars)
	})
}

func handlerNotFound(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		app.NewTemplatingEngine().
			Set("title", app.Locale.T(app.Session.GetLocale(r), "Not found")).
			SetError("title", app.Locale.T(app.Session.GetLocale(r), "Not found")).
			Render(w, r, "404.html")
	})
}

func handlerError(app App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		app.NewTemplatingEngine().
			Set("error", "TODO"). // TODO: show error
			Set("stacktrace", string(debug.Stack())).
			Render(w, r, "error.html")
	})
}

func flash(app App, w http.ResponseWriter, r *http.Request, msg string) {
	_ = app.Session.Flash(w, r, gotuna.NewFlash(msg))
}

func t(app App, r *http.Request, s string) string {
	return app.Locale.T(app.Session.GetLocale(r), s)
}
