package gotdd

import (
	"fmt"
	"net/http"
)

const flashKey = "_flash"

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
