package goprox_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/damnpoet/goprox"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar"))
})

func assertHeaders(t *testing.T, resHeaders http.Header, reqHeaders map[string]string) {
	for name, value := range reqHeaders {
		if actual := strings.Join(resHeaders[name], ", "); actual != value {
			t.Errorf("Invalid header '%s', wanted '%s', got '%s'", name, value, actual)
		}
	}
}

func TestNoConfig(t *testing.T) {
	prox := goprox.New("http://foo.com", goprox.Options{
	// Intentionally left blank.
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	prox.Handler(testHandler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"X-Forwarded-For": "example.com",
	})
}

func TestNonProxyRequest(t *testing.T) {
	prox := goprox.New("http://foo.com", goprox.Options{
		Path: "/proxy",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	prox.Handler(testHandler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"X-Forwarded-For": "",
	})
}

func TestRestrictedToPathRequest(t *testing.T) {
	prox := goprox.New("http://foo.com", goprox.Options{
		Path: "/proxy",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/proxy/foo", nil)

	prox.Handler(testHandler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"X-Forwarded-For": "example.com",
	})
}
