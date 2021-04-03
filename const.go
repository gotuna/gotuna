package gotdd

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/alcalbg/gotdd/static"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const ContentTypeHTML = "text/html; charset=utf-8"

type Options struct {
	Logger         *log.Logger
	Router         *mux.Router
	FS             fs.FS
	Session        *Session
	UserRepository UserRepository
	StaticPrefix   string
	Locale         Locale
}

func OptionsWithDefaults(options Options) Options {
	keyPairs := os.Getenv("APP_KEY")

	if options.Session == nil {
		options.Session = NewSession(sessions.NewCookieStore([]byte(keyPairs)))
	}

	if options.Router == nil {
		options.Router = mux.NewRouter()
	}

	if options.Locale == nil {
		options.Locale = NewLocale(Translations)
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
