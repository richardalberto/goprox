package goprox

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Proxy represents the reverse proxy
type Proxy struct {
	// debug flag will enable/disable debug level logs
	debug bool
	// destination will especify the end host for all proxied requests
	dest *url.URL
	// path will restrict the proxy to the path provided
	path string
}

// New creates a new instance of Proxy
func New(dest *url.URL, options Options) *Proxy {
	return &Proxy{
		dest:  dest,
		debug: options.Debug,
		path:  options.Path,
	}
}

// Handler apply the Proxy specification on the request
func (p *Proxy) Handler(h http.Handler) http.Handler {
	var handler http.Handler

	// proxy handler
	handler = p.handler(h)

	return handler
}

func (p *Proxy) handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.dest.String() == "" || (p.path != "" && !strings.HasPrefix(r.URL.Path, p.path)) {
			h.ServeHTTP(w, r)
			return
		}

		prox := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = p.dest.Scheme
				req.URL.Host = p.dest.Host
				req.URL.Path = strings.Replace(req.URL.Path, p.path, "", 1)
			},
		}

		log.Printf("Forwarding request: %s", r.URL)
		w.Header().Set("X-Forwarded-For", p.dest.Host)

		// TODO: this a workaround to duplicating Access-Control-Allow-Origin header
		// which is not allowed. Fix!
		w.Header().Del("Access-Control-Allow-Origin")

		// send request to destination
		prox.ServeHTTP(w, r)
	})
}
