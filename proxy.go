package goprox

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/rs/cors"

	"gopkg.in/redis.v3"
)

// Proxy represents the reverse proxy
type Proxy struct {
	// debug flag will enable/disable debug level logs
	debug bool
	// cache flag will enable/disable cache
	cache bool
	// destination will especify the end host for all proxied requests
	dest string
	// path will restrict the proxy to the path provided
	path string
	// redis is the instance of the redis client used to cache requests
	redis *redis.Client
	// enableCORS adds cors headers to each request
	enableCORS bool
}

// New creates a new instance of Proxy
func New(dest string, options Options) *Proxy {
	return &Proxy{
		dest:       dest,
		debug:      options.Debug,
		path:       options.Path,
		cache:      options.Cache,
		enableCORS: options.EnableCORS,
	}
}

// Handler apply the Proxy specification on the request
func (p *Proxy) Handler(h http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.dest != "" && strings.HasPrefix(r.URL.Path, p.path) {
			prox := &httputil.ReverseProxy{
				Director: func(req *http.Request) {
					req.URL.Scheme = "http"
					req.URL.Host = p.dest
				},
			}

			log.Printf("Forwarding request: %s", r.URL)
			w.Header().Set("X-Forwarded-For", r.URL.Host)
			prox.ServeHTTP(w, r)
			return
		} else if p.dest == "" {
			log.Printf("Destination URL is empty")
		}

		log.Printf("Non-proxy request: %s", r.URL.String())
		h.ServeHTTP(w, r)
	})

	if p.enableCORS {
		c := cors.New(cors.Options{
			AllowedHeaders: []string{"Authorization", "Accept", "Content-Type"},
			AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		})

		return c.Handler(handler)
	}

	return handler
}
