package goprox_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/damnpoet/goprox"
)

var (
	destURL, _ = url.Parse("http://foo.com")
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
	prox := goprox.New(destURL, goprox.Options{
	// Intentionally left blank.
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	prox.Handler(testHandler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"X-Forwarded-For": "foo.com",
	})
}

func TestNonProxyRequest(t *testing.T) {
	prox := goprox.New(destURL, goprox.Options{
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
	prox := goprox.New(destURL, goprox.Options{
		Path: "/proxy",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/proxy/foo", nil)

	prox.Handler(testHandler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"X-Forwarded-For": "foo.com",
	})
}

func TestCachedRequest(t *testing.T) {
	prox := goprox.New(destURL, goprox.Options{
		Cache: true,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	prox.Handler(testHandler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"X-Forwarded-For": "foo.com",
		"X-Cache":         "SKIP",
	})
}
