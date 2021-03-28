package i18n

import "fmt"

type Translator interface {
	T(s string, p ...interface{}) string
}

func NewTranslator(set map[string]string) Translator {
	return &language{set: set}
}

type language struct {
	set map[string]string
}

// T is short for Translate
func (c language) T(key string, p ...interface{}) string {
	if c.set[key] == "" {
		return "^" + key // mark missing translations
	}

	return fmt.Sprintf(c.set[key], p...)
}

// TODO: date formatting
