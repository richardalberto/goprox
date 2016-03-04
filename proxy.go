package goprox

import "github.com/damnpoet/goprox/rest"

type Proxy struct {
	Username string
	Password string

	Debug bool

	rest *rest.RestClient

	Redis string

	server *Server
}

func New(url, username, password string) *Proxy {
	rest := rest.NewRestClient(username, password, url)

	return &Proxy{
		Username: username,
		Password: password,
		rest:     rest,
		Debug:    true,
		server:   NewServer(rest, true),
	}
}

func (p *Proxy) EnableCache(addr, password string, db int64) error {
	return p.server.EnableCache(addr, password, db)
}

func (p *Proxy) Listen(addr string) {
	p.server.Start(addr)
}
