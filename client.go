// Package gocartel provides a comprehensive API wrapper for the BigCartel API.
//
// gocartel offers an idiomatic way to interact with the BigCartel
// platform in Go applications
package gocartel

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type BigCartelClientHeaders struct {
	UserAgent,
	Accept,
	ContentType,
	Authorization string
}

type BigCartelClient struct {
	BaseURL         string
	Client          *http.Client
	Headers         *BigCartelClientHeaders
	TimeoutDuration time.Duration
}

type ClientOpts struct {
	BaseURL,
	UserAgent,
	BasicAuth string
	TimeoutDuration time.Duration
}

func NewClient(opts ClientOpts) BigCartelClient {
	newClient := &http.Client{}
	headers := &BigCartelClientHeaders{
		UserAgent:     opts.UserAgent,
		Accept:        "application/vnd.api+json",
		ContentType:   "application/vnd.api+json",
		Authorization: opts.BasicAuth,
	}
	if opts.TimeoutDuration == 0 {
		opts.TimeoutDuration = 60 * time.Second
	}
	return BigCartelClient{
		BaseURL:         opts.BaseURL,
		Client:          newClient,
		Headers:         headers,
		TimeoutDuration: opts.TimeoutDuration,
	}
}

func (c BigCartelClient) get(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	return c.Client.Do(req)
}

func (c BigCartelClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.Headers.UserAgent)
	req.Header.Set("Accept", c.Headers.Accept)
	req.Header.Set("Content-Type", c.Headers.ContentType)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.Headers.Authorization))
}
