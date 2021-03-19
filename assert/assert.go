package assert

import (
	"strings"
	"testing"
)

func Equal(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf(`%s: got "%v" want "%v"`, t.Name(), got, want)
	}
}

func Contains(t *testing.T, s, substring string) {
	if !strings.Contains(s, substring) {
		t.Errorf(`%s: string "%s" does not contain "%s"`, t.Name(), s, substring)
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`%s: didn't expect error %v`, t.Name(), err)
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		t.Errorf(`%s: error was expected but didn't get one`, t.Name())
	}
}
