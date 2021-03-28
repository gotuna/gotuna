package i18n

import "fmt"

type Locale interface {
	T(s string, p ...interface{}) string
}

func NewLocale(set map[string]string) Locale {
	return &locale{set: set}
}

type locale struct {
	set map[string]string
}

// T is short for Translate
func (c locale) T(key string, p ...interface{}) string {
	if c.set[key] == "" {
		return "^" + key // mark missing translations
	}

	return fmt.Sprintf(c.set[key], p...)
}

// TODO: date formatting
