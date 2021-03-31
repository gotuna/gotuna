package doubles

import (
	"net/http"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/util"
)

func NewAppStub() http.Handler {
	return app.NewApp(util.Options{
		Logger:         NewLoggerStub(),
		FS:             NewFileSystemStub(nil),
		Session:        session.NewSession(NewGorillaSessionStoreSpy(session.GuestSID)),
		UserRepository: NewUserRepositoryStub(),
	})
}
