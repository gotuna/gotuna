package i18n

import "fmt"

type Locale interface {
	T(language string, s string, p ...interface{}) string
}

func NewLocale(set map[string]map[string]string) Locale {
	return &locale{set}
}

type locale struct {
	set map[string]map[string]string
}

// T is short for Translate
func (c locale) T(language string, key string, p ...interface{}) string {

	if c.set[key][language] == "" {
		return "^" + key // mark missing translations
	}

	return fmt.Sprintf(c.set[key][language], p...)
}

// TODO: date formatting
// TODO: pluralization
