package goprox

import (
	"net/http"
	"regexp"

	"gopkg.in/redis.v3"

	log "github.com/Sirupsen/logrus"
	"github.com/rs/cors"

	"github.com/damnpoet/goprox/rest"
)

type Server struct {
	debug  bool
	client *rest.RestClient

	redis *redis.Client
}

// New creates a new instance of the ded mocked endpoint server
func NewServer(client *rest.RestClient, debug bool) *Server {
	return &Server{
		debug:  debug,
		client: client,
	}
}

func (s *Server) EnableCache(addr, password string, db int64) error {
	s.redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := s.redis.Ping().Result(); err != nil {
		return err
	}

	return nil
}

// Start non-blocking wrapper to ListenAndServe
func (s *Server) Start(url string) {
	s.start(url)
}

// start listening. this will be called in a separated goroutine.
func (s *Server) start(url string) {
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

	log.Infof("Listening on %s", url)
	log.Fatal(http.ListenAndServe(url, handler))
}
