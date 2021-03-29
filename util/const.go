package util

import (
	"io/fs"
	"log"
	"strings"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/models"
	"github.com/alcalbg/gotdd/session"
)

const ContentTypeHTML = "text/html; charset=utf-8"

var GuestRoutes = map[string]string{
	"/login":    "login",
	"/register": "register",
}

type Options struct {
	Logger         *log.Logger
	FS             fs.FS
	Session        *session.Session
	UserRepository models.UserRepository
	StaticPrefix   string
	Locale         i18n.Locale
	GuestRoutes    map[string]string
}

func OptionsWithDefaults(options Options) Options {

	if options.GuestRoutes == nil {
		options.GuestRoutes = GuestRoutes
	}

	if options.Locale == nil {
		options.Locale = i18n.NewLocale(i18n.En)
	}

	// path prefix for static files
	// e.g. "/public" or "http://cdn.example.com/assets"
	options.StaticPrefix = strings.TrimRight(options.StaticPrefix, "/")

	return options
}
