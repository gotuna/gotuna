package gotuna

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "app_session"

type Session struct {
	Store sessions.Store
}

// NewSession returns new session with requested store
func NewSession(store sessions.Store) *Session {
	if store == nil {
		panic("Must supply a valid session store")
	}

	return &Session{Store: store}
}

func (s Session) Put(w http.ResponseWriter, r *http.Request, key string, value string) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("cannot get session from the store")
	}

	// TODO: lock needed?
	session.Values[key] = value

	return s.Store.Save(r, w, session)
}

func (s Session) Get(r *http.Request, key string) (string, error) {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return "", errors.New("cannot get session from the store")
	}

	// TODO: lock needed?
	value, ok := session.Values[key].(string)
	if !ok {
		return "", fmt.Errorf("session holds no value for key %s", key)
	}

	return value, nil
}

func (s Session) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("cannot get session from the store")
	}

	delete(session.Values, key)

	return s.Store.Save(r, w, session)
}

func (s Session) Destroy(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("cannot get session from the store")
	}

	delete(session.Values, UserIDKey)
	session.Options.MaxAge = -1

	return s.Store.Save(r, w, session)
}

func TypeFromString(raw string, t interface{}) error {
	return json.Unmarshal([]byte(raw), &t)
}

func TypeToString(t interface{}) (string, error) {
	b, err := json.Marshal(t)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
