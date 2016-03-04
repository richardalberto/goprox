package goprox

import "github.com/damnpoet/goprox/rest"

type cachedResponse struct {
	Status  int
	RawText string
}

func newCachedResponse(resp rest.Response) *cachedResponse {
	return &cachedResponse{
		Status:  resp.Status(),
		RawText: resp.RawText(),
	}
}
