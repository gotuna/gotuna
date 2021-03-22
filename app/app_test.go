package app_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/gorilla/sessions"
)

func newStubLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func newStubSessionStore(r *http.Request, userSID string) sessions.Store {

	userSession := sessions.NewSession(nil, session.SessionName)
	userSession.Values[session.UserSIDKey] = userSID

	return &stubSessionStore{userSession}
}

type stubSessionStore struct {
	userSession *sessions.Session
}

func (session *stubSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return session.userSession, nil
}

func (session *stubSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return session.userSession, nil
}

func (session *stubSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

func TestRoutes(t *testing.T) {
	routes := []struct {
		route  string
		method string
		status int
	}{
		{"/", http.MethodGet, http.StatusFound},
		{"/", http.MethodPost, http.StatusMethodNotAllowed},
		{"/invalid", http.MethodGet, http.StatusNotFound},
		{"/login", http.MethodGet, http.StatusOK},
	}

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s", r.route), func(t *testing.T) {

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			srv := app.NewServer(newStubLogger(), newStubSessionStore(request, ""))

			srv.Router.ServeHTTP(response, request)

			assert.Equal(t, response.Code, r.status)
		})
	}
}

func TestGuestIsRedirectedToLogin(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	srv := app.NewServer(newStubLogger(), newStubSessionStore(request, ""))

	srv.Router.ServeHTTP(response, request)

	assert.Redirects(t, response, "/login", http.StatusFound)
}

func TestLoggedInUserCanSeeHome(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	session := newStubSessionStore(request, "123")

	srv := app.NewServer(newStubLogger(), session)

	srv.Router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
}
