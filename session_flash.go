package gotdd

import (
	"encoding/json"
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

	err = typeFromString(raw, &messages)
	if err != nil {
		return fmt.Errorf("cannot reconstruct type from json string %v", err)
	}

	messages = append(messages, flashMessage)

	str, err := typeToString(messages)
	if err != nil {
		return fmt.Errorf("cannot convert type to json string %v", err)
	}

	return s.Put(w, r, flashKey, str)
}

func (s Session) Flashes(w http.ResponseWriter, r *http.Request) ([]FlashMessage, error) {

	var messages []FlashMessage

	raw, err := s.Get(r, flashKey)
	if err != nil {
		return messages, nil
	}

	err = typeFromString(raw, &messages)
	if err != nil {
		return messages, fmt.Errorf("cannot get a type from json string %v", err)
	}

	s.Delete(w, r, flashKey)

	return messages, nil
}

func typeFromString(raw string, t interface{}) error {
	return json.Unmarshal([]byte(raw), &t)
}

func typeToString(t interface{}) (string, error) {
	b, err := json.Marshal(t)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
