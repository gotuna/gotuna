package i18n_test

import (
	"testing"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestTranslations(t *testing.T) {

	var translationsStub = map[string]map[string]string{
		"The world": {
			"en-US": "The world",
			"de-DE": "Die welt",
		},
		"Page %d of %d": {
			"en-US": "Page %d of %d",
			"it-IT": "Pagina %d di %d",
		},
		"The %s color has a value of %d": {
			"en-US": "The %s color has a value of %d",
			"es-ES": "El color %s tiene un valor de %d",
		},
	}

	locale := i18n.NewLocale(translationsStub)

	assert.Equal(t, "Die welt", locale.T("de-DE", "The world"))
	assert.Equal(t, "Pagina 2 di 4", locale.T("it-IT", "Page %d of %d", 2, 4))
	assert.Equal(t, "El color Rojo tiene un valor de 10", locale.T("es-ES", "The %s color has a value of %d", "Rojo", 10))
	assert.Equal(t, "^Unknown string", locale.T("en-US", "Unknown string"))
	assert.Equal(t, "^Unknown string", locale.T("ru-RU", "Unknown string"))
}
