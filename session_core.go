package gotuna

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// ErrCannotGetSession is thrown when we cannot retrieve a valid session from the store
	ErrCannotGetSession = constError("cannot get session from the store")
	// ErrNoValueForThisKey is thrown when we cannot get a value for provided key
	ErrNoValueForThisKey = constError("session holds no value for this key")
)

// Session is the main application session store.
type Session struct {
	store sessions.Store
	name  string
}

// NewSession returns a new application session with requested store engine.
func NewSession(store sessions.Store, name string) *Session {
	if store == nil {
		panic("must supply a valid session store")
	}

	if name == "" {
		panic("must supply a valid session name")
	}

	return &Session{
		store: store,
		name:  name,
	}
}

// Put string value in the session for specified key.
func (s Session) Put(w http.ResponseWriter, r *http.Request, key string, value string) error {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		return ErrCannotGetSession
	}

	// TODO: lock needed?
	session.Values[key] = value

	return s.store.Save(r, w, session)
}

// Get string value from the session for specified key.
func (s Session) Get(r *http.Request, key string) (string, error) {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		return "", ErrCannotGetSession
	}

	// TODO: lock needed?
	value, ok := session.Values[key].(string)
	if !ok {
		return "", ErrNoValueForThisKey
	}

	return value, nil
}

// Delete value from the session for key.
func (s Session) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		return ErrCannotGetSession
	}

	delete(session.Values, key)

	return s.store.Save(r, w, session)
}

// Destroy the user session by removing the user key and expiring the cookie.
func (s Session) Destroy(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		return ErrCannotGetSession
	}

	delete(session.Values, UserIDKey)
	session.Options.MaxAge = -1

	return s.store.Save(r, w, session)
}

// TypeToString converts any type t to JSON-encoded string value.
// This is used because app's session can only hold basic string values.
func TypeToString(t interface{}) (string, error) {
	b, err := json.Marshal(t)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// TypeFromString converts JSON-encoded value into the type t.
func TypeFromString(raw string, t interface{}) error {
	return json.Unmarshal([]byte(raw), &t)
}
