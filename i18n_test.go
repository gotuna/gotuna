package gotuna_test

import (
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
)

func TestTranslations(t *testing.T) {

	locale := gotuna.NewLocale(map[string]map[string]string{
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
	})

	assert.Equal(t, "Die welt", locale.T("de-DE", "The world"))
	assert.Equal(t, "Pagina 2 di 4", locale.T("it-IT", "Page %d of %d", 2, 4))
	assert.Equal(t, "El color Rojo tiene un valor de 10", locale.T("es-ES", "The %s color has a value of %d", "Rojo", 10))
	assert.Equal(t, "Unknown string", locale.T("en-US", "Unknown string"))
	assert.Equal(t, "Unknown string", locale.T("ru-RU", "Unknown string"))
}

func TestPluralization(t *testing.T) {

	locale := gotuna.NewLocale(map[string]map[string]string{
		"oranges": {
			"en-US": "There is one orange|There are many oranges",
		},
		"apples": {
			"en-US": "%s apple|%s apples",
		},
	})

	assert.Equal(t, "There is one orange", locale.TP("en-US", "oranges", 1))
	assert.Equal(t, "There are many oranges", locale.TP("en-US", "oranges", 5))
	assert.Equal(t, "green apple", locale.TP("en-US", "apples", 1, "green"))
	assert.Equal(t, "red apples", locale.TP("en-US", "apples", 22, "red"))
}
