package gotdd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const GuestSID = ""
const UserSIDKey = "_user_sid"
const flashKey = "_flash"
const sessionName = "app_session"

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

func (s Session) Get(r *http.Request, key string) (string, error) {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return "", errors.New("cannot get session from the store")
	}

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

func (s Session) SetUserSID(w http.ResponseWriter, r *http.Request, sid string) error {
	return s.Put(w, r, UserSIDKey, sid)
}

func (s Session) GetUserSID(r *http.Request) (string, error) {
	sid, err := s.Get(r, UserSIDKey)
	if err != nil || sid == GuestSID {
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

	var messages []FlashMessage

	raw, err := s.Get(r, flashKey)
	if err != nil {
		raw = "[]"
	}

	err = json.Unmarshal([]byte(raw), &messages)
	if err != nil {
		return err
	}

	messages = append(messages, flashMessage)

	rawbytes, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	return s.Put(w, r, flashKey, string(rawbytes))
}

func (s Session) Flashes(w http.ResponseWriter, r *http.Request) ([]FlashMessage, error) {

	var messages []FlashMessage

	raw, err := s.Get(r, flashKey)
	if err != nil {
		return messages, nil
	}

	err = json.Unmarshal([]byte(raw), &messages)
	if err != nil {
		return messages, err
	}

	s.Delete(w, r, flashKey)

	return messages, nil
}
