package gotdd

import (
	"io"
	"io/fs"
	"log"

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
		app.Logger = log.New(io.Discard, "", 0)
	}

	if app.Locale == nil {
		app.Locale = NewLocale(Translations)
	}

	return app
}
