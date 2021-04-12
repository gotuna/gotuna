package gotuna

import (
	"fmt"
	"net/http"
)

const flashKey = "_flash"

// FlashMessage represents a flash message.
type FlashMessage struct {
	Message   string
	Kind      string
	AutoClose bool
}

// NewFlash is a constructor for new flash message.
func NewFlash(message string) FlashMessage {
	return FlashMessage{
		Message:   message,
		Kind:      "success",
		AutoClose: true,
	}
}

// Flash will store flash message into the session.
// These messages are primarily used to inform the user during a subsequesnt
// request about some status update.
func (s Session) Flash(w http.ResponseWriter, r *http.Request, flashMessage FlashMessage) error {

	var messages []FlashMessage

	raw, err := s.Get(r, flashKey)
	if err != nil {
		raw = "[]"
	}

	err = TypeFromString(raw, &messages)
	if err != nil {
		return fmt.Errorf("cannot reconstruct type from json string %v", err)
	}

	messages = append(messages, flashMessage)

	str, err := TypeToString(messages)
	if err != nil {
		return fmt.Errorf("cannot convert type to json string %v", err)
	}

	return s.Put(w, r, flashKey, str)
}

// Flashes returns all messages from the session and removes them.
func (s Session) Flashes(w http.ResponseWriter, r *http.Request) []FlashMessage {

	var messages []FlashMessage

	raw, err := s.Get(r, flashKey)
	if err != nil {
		return messages
	}

	TypeFromString(raw, &messages)

	s.Delete(w, r, flashKey)

	return messages
}
