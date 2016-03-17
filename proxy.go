package goprox

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/lox/httpcache"
	"github.com/rs/cors"
)

// Proxy represents the reverse proxy
type Proxy struct {
	// debug flag will enable/disable debug level logs
	debug bool
	// cache flag will enable/disable cache
	cache bool
	// destination will especify the end host for all proxied requests
	dest *url.URL
	// path will restrict the proxy to the path provided
	path string
	// cors adds cors headers to each request
	cors bool
}

// New creates a new instance of Proxy
func New(dest *url.URL, options Options) *Proxy {
	return &Proxy{
		dest:  dest,
		debug: options.Debug,
		path:  options.Path,
		cache: options.Cache,
		cors:  options.CORS,
	}
}

// Handler apply the Proxy specification on the request
func (p *Proxy) Handler(h http.Handler) http.Handler {
	var handler http.Handler

	// proxy handler
	handler = p.handler(h)

	// append httpcache middleware
	if p.cache {
		handler = p.cacheHandler(handler)
	}

	// append cors middleware
	if p.cors {
		handler = p.corsHandler(handler)
	}

	return handler
}

func (p *Proxy) handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.dest.String() == "" || (p.path != "" && !strings.HasPrefix(r.URL.Path, p.path)) {
			log.Printf("Missing destination")
			h.ServeHTTP(w, r)
		} else {
			prox := &httputil.ReverseProxy{
				Director: func(req *http.Request) {
					req.URL.Scheme = p.dest.Scheme
					req.URL.Host = p.dest.Host
				},
			}

			log.Printf("Forwarding request: %s", r.URL)
			w.Header().Set("X-Forwarded-For", p.dest.Host)

			// TODO: this a workaround to duplicating Access-Control-Allow-Origin header
			// which is not allowed. Fix!
			w.Header().Del("Access-Control-Allow-Origin")

			// send request to destination
			prox.ServeHTTP(w, r)
		}
	})
}

func (p *Proxy) corsHandler(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization", "Accept", "Content-Type"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		Debug:          p.debug,
	})

	return c.Handler(h)
}

func (p *Proxy) cacheHandler(h http.Handler) http.Handler {
	c := httpcache.NewHandler(httpcache.NewMemoryCache(), h)
	c.Shared = true

	return c
}
