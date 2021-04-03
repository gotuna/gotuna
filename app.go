package gotdd

import (
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
