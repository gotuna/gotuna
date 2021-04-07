package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/cmd/main/i18n"
	"github.com/alcalbg/gotdd/cmd/main/static"
	"github.com/alcalbg/gotdd/cmd/main/views"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {

	port := ":8888"
	keyPairs := os.Getenv("APP_KEY")

	app := MakeApp(gotdd.App{
		Logger:         log.New(os.Stdout, "", 0),
		UserRepository: doubles.NewUserRepositoryStub(),
		Session:        gotdd.NewSession(sessions.NewCookieStore([]byte(keyPairs))),
		Static:         static.EmbededStatic,
		StaticPrefix:   "",
		Views:          views.EmbededViews,
		Locale:         gotdd.NewLocale(i18n.Translations),
	})

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

func MakeApp(app gotdd.App) gotdd.App {

	if app.Logger == nil {
		app.Logger = log.New(io.Discard, "", 0)
	}

	if app.Locale == nil {
		app.Locale = gotdd.NewLocale(map[string]map[string]string{})
	}

	app.ViewHelpers = template.FuncMap{
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	app.Router = mux.NewRouter()
	app.Router.NotFoundHandler = handlerNotFound(app)

	// middlewares for all routes
	app.Router.Handle("/error", handlerError(app)).Methods(http.MethodGet, http.MethodPost)
	app.Router.Use(app.Recoverer("/error"))
	app.Router.Use(app.Logging())
	// TODO: csrf middleware
	app.Router.Methods(http.MethodOptions)
	app.Router.Use(app.Cors())

	// logged in user
	user := app.Router.NewRoute().Subrouter()
	user.Use(app.Authenticate("/login"))
	user.Use(app.StoreUserToContext())
	user.Handle("/", handlerHome(app)).Methods(http.MethodGet)
	user.Handle("/profile", handlerProfile(app)).Methods(http.MethodGet, http.MethodPost)
	user.Handle("/logout", handlerLogout(app)).Methods(http.MethodPost)
	user.Handle("/setlocale/{locale}", handlerChangeLocale(app)).Methods(http.MethodGet, http.MethodPost)

	// guests
	auth := app.Router.NewRoute().Subrouter()
	auth.Use(app.RedirectIfAuthenticated("/"))
	auth.Handle("/login", handlerLogin(app)).Methods(http.MethodGet, http.MethodPost)

	// path prefix for static files
	// e.g. "/public" or "http://cdn.example.com/assets"
	app.StaticPrefix = strings.TrimRight(app.StaticPrefix, "/")

	// serve files from the static directory
	app.Router.PathPrefix(app.StaticPrefix).
		Handler(http.StripPrefix(app.StaticPrefix, serveFiles(app))).
		Methods(http.MethodGet)

	return app
}

func serveFiles(app gotdd.App) http.Handler {
	fs := http.FS(app.Static)
	fileapp := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			handlerNotFound(app).ServeHTTP(w, r)
			return
		}
		stat, _ := f.Stat()
		if stat.IsDir() {
			handlerNotFound(app).ServeHTTP(w, r)
			return
		}

		// TODO: ModTime doesn't work for embed?
		//w.Header().Set("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", "31536000"))
		fileapp.ServeHTTP(w, r)
	})
}

func handlerHome(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewNativeTemplatingEngine().
			Set("message", app.Locale.T(app.Session.GetUserLocale(r), "Home")).
			Render(w, r, "app.html", "home.html")
	})
}

func handlerLogin(app gotdd.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl := app.NewNativeTemplatingEngine()

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

		if err := app.Session.Flash(w, r, gotdd.NewFlash(app.Locale.T(app.Session.GetUserLocale(r), "Welcome"))); err != nil {
			app.Logger.Printf("%s %s %s %v", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, err)
			handlerError(app).ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func handlerLogout(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Session.Destroy(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func handlerProfile(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewNativeTemplatingEngine().
			Render(w, r, "app.html", "profile.html")
	})
}

func handlerChangeLocale(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		app.Session.SetUserLocale(w, r, vars["locale"])
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func handlerNotFound(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		app.NewNativeTemplatingEngine().
			Set("title", app.Locale.T(app.Session.GetUserLocale(r), "Not found")).
			Render(w, r, "app.html", "4xx.html")
	})
}

func handlerError(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		app.NewNativeTemplatingEngine().
			Set("error", "TODO"). // TODO: show error
			Set("stacktrace", string(debug.Stack())).
			Render(w, r, "app.html", "error.html")
	})
}
