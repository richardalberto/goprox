package goprox

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func (s *Server) SpadeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Infof("Received request: %s", r.URL)
		s.get(w, r)
	}
	// TODO: support other methods
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// make a copy of client
	client := *s.client
	if auth := r.Header.Get("Authorization"); auth != "" {
		client.Header.Set("Authorization", auth)
	}

	resp, err := client.Get(r.URL.Path)
	if err != nil {
		log.Errorf("An error ocurred %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		if buf, err := NewError(err).JSON(); err != nil {
			w.Write(buf)
		}
		return
	}

	w.WriteHeader(resp.Status())
	w.Write([]byte(resp.RawText()))
}
