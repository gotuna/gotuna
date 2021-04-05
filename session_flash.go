package gotdd

import (
	"encoding/json"
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
