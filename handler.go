package goprox

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

func (s *Server) SpadeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Infof("Received GET Request: %s", r.URL)
		s.get(w, r)
	case "PUT":
		log.Infof("Received PUT Request: %s", r.URL)
		s.put(w, r)
	case "POST":
		log.Infof("Received POST Request: %s", r.URL)
		s.post(w, r)
	default:
		log.Errorf("Recived Request with Invalid Method: %s", r.Method)
	}
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// try to use cache if available
	if s.redis != nil {
		if cachedResp, err := s.redis.Get(r.URL.Path).Bytes(); err == nil {
			var resp cachedResponse
			if err = json.Unmarshal(cachedResp, &resp); err == nil {
				log.Infof("Using cached version of %s", r.URL.Path)
				w.WriteHeader(resp.Status)
				w.Write([]byte(resp.RawText))
				return
			}
		}
		log.Infof("Cached version of %s not found, doing the actual request", r.URL.Path)
	}

	// make a copy of client
	client := *s.client
	if auth := r.Header.Get("Authorization"); auth != "" {
		client.Header.Set("Authorization", auth)
	}

	// do request
	resp, err := client.Get(r.URL.Path)
	if err != nil {
		log.Errorf("An error ocurred %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		if buf, err := NewError(err).JSON(); err != nil {
			w.Write(buf)
		}
		return
	}

	// cache response
	if s.redis != nil {
		if serialized, err := json.Marshal(newCachedResponse(*resp)); err == nil {
			if err := s.redis.Set(r.URL.Path, serialized, time.Second*5).Err(); err != nil {
				log.Errorf("An error ocurred while trying to cache the response for GET: %s in cache, %s", r.URL.Path, err)
			}
		}
	}

	// write response
	w.WriteHeader(resp.Status())
	w.Write([]byte(resp.RawText()))
}

func (s *Server) put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// make a copy of client
	client := *s.client
	if auth := r.Header.Get("Authorization"); auth != "" {
		client.Header.Set("Authorization", auth)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error while reading POST body for %s, %s", r.URL.Path, err)
		return
	}

	resp, err := client.Put(r.URL.Path, bytes.NewBuffer(b))
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

func (s *Server) post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// make a copy of client
	client := *s.client
	if auth := r.Header.Get("Authorization"); auth != "" {
		client.Header.Set("Authorization", auth)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error while reading POST body for %s, %s", r.URL.Path, err)
		return
	}

	resp, err := client.Post(r.URL.Path, bytes.NewBuffer(b))
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
