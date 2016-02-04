package rest

import (
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

type RestClient struct {
	client *Client
	url    string
}

func NewRestClient(username, password, path string) *RestClient {
	c := &Client{}
	c.UserInfo = url.UserPassword(username, password)

	return &RestClient{
		client: c,
		url:    path,
	}
}

func (c *RestClient) Get(path string) (*Response, error) {
	uri := fmt.Sprintf("%s/%s", c.url, path)

	resp, err := c.client.Get(uri)
	if err != nil {
		log.Error(err)
	}

	return resp, err
}
