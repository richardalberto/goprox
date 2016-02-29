package goprox

import (
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/rs/cors"

	"github.com/damnpoet/goprox/rest"
)

type Server struct {
	url    string
	debug  bool
	client *rest.RestClient
}

// New creates a new instance of the ded mocked endpoint server
func NewServer(url string, client *rest.RestClient, debug bool) *Server {
	return &Server{
		url:    url,
		debug:  debug,
		client: client,
	}
}

// Start non-blocking wrapper to ListenAndServe
func (s *Server) Start() {
	s.start()
}

// start listening. this will be called in a separated goroutine.
func (s *Server) start() {
	regex, _ := regexp.Compile("/*")

	regexHandler := NewRegexpHandler()
	regexHandler.HandleFunc(regex, s.SpadeHandler)

	// add cors
	// TODO: make this optional
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization", "Accept", "Content-Type"},
		Debug:          s.debug,
		AllowedMethods: []string{"GET", "PUT"},
	})
	handler := c.Handler(regexHandler)

	log.Infof("Listening on %s", s.url)
	log.Fatal(http.ListenAndServe(s.url, handler))
}
