package gotuna_test

import (
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
)

func TestErrorConstants(t *testing.T) {
	assert.Equal(t, "this field is required", gotuna.ErrRequiredField.Error())

	var err error = gotuna.ErrRequiredField
	assert.Equal(t, err, gotuna.ErrRequiredField)
}
