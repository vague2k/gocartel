package gocartel

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type clientHeaders struct {
	userAgent,
	accept,
	contentType,
	authorization string
}

type bigCartelClient struct {
	baseURL         string
	client          *http.Client
	headers         *clientHeaders
	timeoutDuration time.Duration
}

type ClientOpts struct {
	BaseURL,
	UserAgent,
	BasicAuth string
	TimeoutDuration time.Duration
}

func NewClient(opts ClientOpts) bigCartelClient {
	newClient := &http.Client{}
	headers := &clientHeaders{
		userAgent:     opts.UserAgent,
		accept:        "application/vnd.api+json",
		contentType:   "application/vnd.api+json",
		authorization: opts.BasicAuth,
	}
	if opts.TimeoutDuration == 0 {
		opts.TimeoutDuration = 60 * time.Second
	}
	return bigCartelClient{
		baseURL:         opts.BaseURL,
		client:          newClient,
		headers:         headers,
		timeoutDuration: opts.TimeoutDuration,
	}
}

func (c bigCartelClient) get(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
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
