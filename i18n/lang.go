package i18n

import "fmt"

type Translator interface {
	Translate(s string, p ...interface{}) string
}

func NewTranslator(set map[string]string) Translator {
	return &language{set: set}
}

type language struct {
	set map[string]string
}

func (c language) Translate(key string, p ...interface{}) string {
	if c.set[key] == "" {
		return "^" + key // mark missing translations
	}

	return fmt.Sprintf(c.set[key], p...)
}
