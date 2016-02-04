package goprox

import "github.com/damnpoet/goprox/rest"

type Proxy struct {
	Username string
	Password string

	Debug bool

	rest *rest.RestClient
}

func New(url, username, password string) *Proxy {
	rest := rest.NewRestClient(username, password, url)

	return &Proxy{
		Username: username,
		Password: password,
		rest:     rest,
		Debug:    true,
	}
}

func (p *Proxy) Listen(addr string) {
	server := NewServer(addr, p.rest, p.Debug)
	server.Start()
}
