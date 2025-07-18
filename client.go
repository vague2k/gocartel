package gocartel

import (
	"fmt"
	"net/http"
)

type clientHeaders struct {
	userAgent,
	accept,
	contentType,
	authorization string
}

type bigCartelClient struct {
	baseURL string
	client  *http.Client
	headers *clientHeaders
}

type ClientOpts struct {
	BaseURL,
	UserAgent,
	BasicAuth string
}

func NewClient(opts ClientOpts) bigCartelClient {
	newClient := &http.Client{}
	headers := &clientHeaders{
		userAgent:     opts.UserAgent,
		accept:        "application/vnd.api+json",
		contentType:   "application/vnd.api+json",
		authorization: opts.BasicAuth,
	}
	return bigCartelClient{
		baseURL: opts.BaseURL,
		client:  newClient,
		headers: headers,
	}
}

func (c bigCartelClient) get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	return c.client.Do(req)
}

func (c bigCartelClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.headers.userAgent)
	req.Header.Set("Accept", c.headers.accept)
	req.Header.Set("Content-Type", c.headers.contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.headers.authorization))
}
