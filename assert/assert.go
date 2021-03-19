package assert

import (
	"testing"
)

func Equal(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf("%s: got %v want %v", t.Name(), got, want)
	}
}
