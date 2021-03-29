package util

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/models"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/static"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const ContentTypeHTML = "text/html; charset=utf-8"

var GuestRoutes = map[string]string{
	"/login":    "login",
	"/register": "register",
}

type Options struct {
	Logger         *log.Logger
	Router         *mux.Router
	FS             fs.FS
	Session        *session.Session
	UserRepository models.UserRepository
	StaticPrefix   string
	Locale         i18n.Locale
	GuestRoutes    map[string]string
}

func OptionsWithDefaults(options Options) Options {
	keyPairs := os.Getenv("APP_KEY")

	if options.Session == nil {
		options.Session = session.NewSession(sessions.NewCookieStore([]byte(keyPairs)))
	}

	if options.Router == nil {
		options.Router = mux.NewRouter()
	}

	if options.GuestRoutes == nil {
		options.GuestRoutes = GuestRoutes
	}

	if options.Locale == nil {
		options.Locale = i18n.NewLocale(i18n.En)
	}

	if options.FS == nil {
		options.FS = static.EmbededStatic
	}

	if options.Logger == nil {
		options.Logger = log.New(os.Stdout, "", 0)
	}

	// path prefix for static files
	// e.g. "/public" or "http://cdn.example.com/assets"
	options.StaticPrefix = strings.TrimRight(options.StaticPrefix, "/")

	return options
}
