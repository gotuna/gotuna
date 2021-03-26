package i18n_test

import (
	"testing"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestTranslations(t *testing.T) {

	var fakeSet = map[string]string{
		"The world":                      "Die welt",
		"Page %d of %d":                  "Pagina %d di %d",
		"The %s color has a value of %d": "El color %s tiene un valor de %d",
	}

	translator := i18n.NewTranslator(fakeSet)

	assert.Equal(t, "Die welt", translator.Translate("The world"))
	assert.Equal(t, "Pagina 2 di 4", translator.Translate("Page %d of %d", 2, 4))
	assert.Equal(t, "El color Rojo tiene un valor de 10", translator.Translate("The %s color has a value of %d", "Rojo", 10))
	assert.Equal(t, "^Unknown string", translator.Translate("Unknown string"))
}
