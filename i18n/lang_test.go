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

	locale := i18n.NewLocale(fakeSet)

	assert.Equal(t, "Die welt", locale.T("The world"))
	assert.Equal(t, "Pagina 2 di 4", locale.T("Page %d of %d", 2, 4))
	assert.Equal(t, "El color Rojo tiene un valor de 10", locale.T("The %s color has a value of %d", "Rojo", 10))
	assert.Equal(t, "^Unknown string", locale.T("Unknown string"))
}
