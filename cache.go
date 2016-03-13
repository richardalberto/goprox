package goprox

import "github.com/damnpoet/goprox/rest"

type cachedResponse struct {
	Status  int
	RawText string
	Auth    string
}

func newCachedResponse(resp rest.Response, auth string) *cachedResponse {
	return &cachedResponse{
		Status:  resp.Status(),
		RawText: resp.RawText(),
		Auth:    auth,
	}
}
