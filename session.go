package gotdd

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const GuestSID = ""
const UserSIDKey = "_user_sid"
const flashKey = "_flash"
const sessionName = "app_session"

func init() {
	gob.Register([]FlashMessage{})
}

type Session struct {
	Store sessions.Store
}

type FlashMessage struct {
	Message   string
	Kind      string
	AutoClose bool
}

func NewFlash(message string) FlashMessage {
	return FlashMessage{
		Message:   message,
		Kind:      "success",
		AutoClose: true,
	}
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

	session.Values[key] = value

	return s.Store.Save(r, w, session)
}

func (s Session) Get(w http.ResponseWriter, r *http.Request, key string) (string, error) {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return "", errors.New("cannot get session from the store")
	}

	sid, ok := session.Values[key].(string)
	if !ok {
		return "", fmt.Errorf("session holds no value for key %s", key)
	}

	return sid, nil
}

func (s Session) SetUserSID(w http.ResponseWriter, r *http.Request, sid string) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("cannot get session from the store")
	}

	session.Values[UserSIDKey] = sid

	return s.Store.Save(r, w, session)
}

func (s Session) GetUserSID(r *http.Request) (string, error) {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return GuestSID, errors.New("cannot get session from the store")
	}

	sid, ok := session.Values[UserSIDKey].(string)
	if !ok || sid == GuestSID {
		return GuestSID, errors.New("No user in the session")
	}

	return sid, nil
}

func (s Session) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("cannot get session from the store")
	}

	delete(session.Values, UserSIDKey)
	session.Options.MaxAge = -1

	return s.Store.Save(r, w, session)
}

func (s Session) IsGuest(r *http.Request) bool {
	sid, err := s.GetUserSID(r)
	if err != nil {
		return true
	}

	if sid == GuestSID {
		return true
	}

	return false
}

func (s Session) Flash(w http.ResponseWriter, r *http.Request, flashMessage FlashMessage) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("cannot get session from the store")
	}

	var flashes []FlashMessage

	if v, ok := session.Values[flashKey]; ok {
		flashes = v.([]FlashMessage)
	}
	session.Values[flashKey] = append(flashes, flashMessage)

	return s.Store.Save(r, w, session)
}

func (s Session) Flashes(w http.ResponseWriter, r *http.Request) ([]FlashMessage, error) {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return nil, errors.New("cannot get session from the store")
	}

	messages, ok := session.Values[flashKey].([]FlashMessage)
	if !ok {
		messages = []FlashMessage{}
	}

	delete(session.Values, flashKey)

	return messages, s.Store.Save(r, w, session)
}
