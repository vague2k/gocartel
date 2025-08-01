// Package gocartel provides a comprehensive API wrapper for the BigCartel API.
//
// gocartel offers an idiomatic way to interact with the BigCartel
// platform in Go applications
package gocartel

import (
	"context"
	"fmt"
	"io"
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
	BaseURL string
	Client  *http.Client
	Headers *BigCartelClientHeaders
}

type ClientOpts struct {
	UserAgent,
	BasicAuth string
	HTTPClient *http.Client // 60 seconds is the default timeout duration if none was given
}

func NewClient(opts ClientOpts) BigCartelClient {
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{}
	}
	if opts.HTTPClient.Timeout == 0 {
		opts.HTTPClient.Timeout = 60 * time.Second
	}
	headers := &BigCartelClientHeaders{
		UserAgent:     opts.UserAgent,
		Accept:        "application/vnd.api+json",
		ContentType:   "application/vnd.api+json",
		Authorization: opts.BasicAuth,
	}
	return BigCartelClient{
		BaseURL: "https://api.bigcartel.com/v1",
		Client:  opts.HTTPClient,
		Headers: headers,
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

func (c BigCartelClient) post(ctx context.Context, endpoint string, data io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+endpoint, data)
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
