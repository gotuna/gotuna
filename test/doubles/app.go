package doubles

import (
	"net/http"

	"github.com/alcalbg/gotdd"
)

func NewAppStub() http.Handler {
	return gotdd.NewApp(gotdd.Options{
		Logger:         NewLoggerStub(),
		FS:             NewFileSystemStub(nil),
		Session:        gotdd.NewSession(NewGorillaSessionStoreSpy(gotdd.GuestSID)),
		UserRepository: NewUserRepositoryStub(),
	})
}
