package lang_test

import (
	"testing"

	"github.com/alcalbg/gotdd/lang"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestTranslations(t *testing.T) {

	var fakeSet = map[string]string{
		"The world":                      "Die welt",
		"Page %d of %d":                  "Pagina %d di %d",
		"The %s color has a value of %d": "El color %s tiene un valor de %d",
	}

	Lang := lang.NewTranslator(fakeSet)

	assert.Equal(t, "Die welt", Lang.T("The world"))
	assert.Equal(t, "Pagina 2 di 4", Lang.T("Page %d of %d", 2, 4))
	assert.Equal(t, "El color Rojo tiene un valor de 10", Lang.T("The %s color has a value of %d", "Rojo", 10))
	assert.Equal(t, "^Unknown string", Lang.T("Unknown string"))
}
