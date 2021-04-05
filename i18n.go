package gotdd

import (
	"fmt"
	"strings"
)

type Locale interface {
	T(language string, s string, p ...interface{}) string
	TP(language string, s string, n int, p ...interface{}) string
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
		return key
	}

	return fmt.Sprintf(c.set[key][language], p...)
}

// TP is short for TranslatePlural
func (c locale) TP(language string, key string, n int, p ...interface{}) string {

	if c.set[key][language] == "" {
		return key
	}

	s := c.set[key][language]
	split := strings.Split(s, "|")

	if n > 1 && len(split) > 1 {
		return fmt.Sprintf(split[1], p...)
	}

	return fmt.Sprintf(split[0], p...)
}

// TODO: date formatting
