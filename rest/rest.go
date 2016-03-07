package rest

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

type RestClient struct {
	client *Client
	url    string
	Header http.Header

	userInfo *url.Userinfo
}

func NewRestClient(username, password, path string) *RestClient {
	return &RestClient{
		client:   &Client{},
		url:      path,
		userInfo: url.UserPassword(username, password),
		Header:   http.Header{},
	}
}

func (c *RestClient) Get(path string) (*Response, error) {
	uri := fmt.Sprintf("%s/%s", c.url, path)

	// make a copy of client
	client := *c.client

	// add authorization header if there's no raw authorization header already.
	if auth := c.Header.Get("Authorization"); auth != "" {
		client.Header = &c.Header
	} else {
		client.UserInfo = c.userInfo
	}

	resp, err := client.Get(uri)
	if err != nil {
		log.Error(err)
	}

	return resp, err
}

func (c *RestClient) Put(path string, payload interface{}) (*Response, error) {
	uri := fmt.Sprintf("%s/%s", c.url, path)

	// make a copy of client
	client := *c.client

	// add authorization header if there's no raw authorization header already.
	if auth := c.Header.Get("Authorization"); auth != "" {
		client.Header = &c.Header
	} else {
		client.UserInfo = c.userInfo
	}

	resp, err := client.Put(uri, payload)
	if err != nil {
		log.Error(err)
	}

	return resp, err
}

func (c *RestClient) Post(path string, payload interface{}) (*Response, error) {
	uri := fmt.Sprintf("%s/%s", c.url, path)

	// make a copy of client
	client := *c.client

	// add authorization header if there's no raw authorization header already.
	if auth := c.Header.Get("Authorization"); auth != "" {
		client.Header = &c.Header
	} else {
		client.UserInfo = c.userInfo
	}

	resp, err := client.Post(uri, payload)
	if err != nil {
		log.Error(err)
	}

	return resp, err
}

func (c *RestClient) Delete(path string) (*Response, error) {
	uri := fmt.Sprintf("%s/%s", c.url, path)

	// make a copy of client
	client := *c.client

	// add authorization header if there's no raw authorization header already.
	if auth := c.Header.Get("Authorization"); auth != "" {
		client.Header = &c.Header
	} else {
		client.UserInfo = c.userInfo
	}

	resp, err := client.Delete(uri)
	if err != nil {
		log.Error(err)
	}

	return resp, err
}
