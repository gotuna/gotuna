package assert

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func Equal(t *testing.T, a, b interface{}) {
	t.Helper()

	if a != b {
		t.Errorf(`%s: "%v" and "%v" should be equal`, t.Name(), a, b)
	}
}

func Greater(t *testing.T, a, b int) {
	t.Helper()

	if !(a > b) {
		t.Errorf(`%s: %d should be greater than %d`, t.Name(), a, b)
	}
}

func Contains(t *testing.T, s, substring string) {
	t.Helper()

	if !strings.Contains(s, substring) {
		t.Errorf(`%s: string "%s" does not contain "%s"`, t.Name(), s, substring)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf(`%s: didn't expect error %v`, t.Name(), err)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf(`%s: error was expected but didn't get one`, t.Name())
	}
}

func Redirects(t *testing.T, r *httptest.ResponseRecorder, url string, code int) {
	t.Helper()

	location, err := r.Result().Location()

	if location == nil || err != nil {
		t.Errorf("response should redirect to %s got status %d", url, r.Code)
		return
	}

	if location.String() != url {
		t.Errorf("response should redirect to %s got status %d", url, r.Code)
	}

	if r.Code != code {
		t.Errorf("response code should be %d got %d", r.Code, code)
	}
}
