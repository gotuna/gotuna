package lang

import "fmt"

type Translator interface {
	T(s string, p ...interface{}) string
}

// global var
var Lang *language

// default language is English
func init() {
	InitTranslator(En)
}

func InitTranslator(set map[string]string) {
	Lang = &language{set: set}
}

type language struct {
	set map[string]string
}

func (c language) T(key string, p ...interface{}) string {
	if c.set[key] == "" {
		// mark missing translations
		return "^" + key
	}

	return fmt.Sprintf(c.set[key], p...)
}
