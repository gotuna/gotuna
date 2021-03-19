package assert

import (
	"testing"
)

func Equal(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
