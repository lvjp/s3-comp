package client

import (
	"context"
	"net/http"
	"net/url"
)

type Client struct {
	config Config
}

func New(cfg Config) (*Client, error) {
	c := &Client{
		config: cfg,
	}

	if err := c.config.SetDefaults(); err != nil {
		return nil, err
	}

	return c, nil
}

func newRequest() *http.Request {
	return &http.Request{
		URL:    new(url.URL),
		Header: http.Header{},

		// The value -1 indicates that the length is unknown.
		ContentLength: -1,
		Body:          http.NoBody,
	}
}

func (c *Client) resolve(ctx context.Context, req *http.Request, input any) error {
	params := EndpointParameters{
		Region:       &c.config.Region,
		Endpoint:     &c.config.Endpoint,
		UsePathStyle: c.config.UsePathStyle,
	}

	if bucketGetter, ok := input.(BucketGetter); ok {
		bucket := bucketGetter.GetBucket()
		params.Bucket = &bucket
	}

	endpoint, err := c.config.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		if _, ok := c.config.EndpointResolver.(*DefaultEndpointResolver); ok {
			return err
		}

		return NewSDKError(ctx, "cannot resolve endpoint: "+err.Error())
	}

	req.URL.Scheme = endpoint.URI.Scheme
	req.URL.Host = endpoint.URI.Host
	req.URL.Path = joinURIPath(endpoint.URI.Path, req.URL.Path)
	req.URL.RawPath = joinURIPath(endpoint.URI.RawPath, req.URL.RawPath)

	for key := range endpoint.Headers {
		req.Header.Set(key, endpoint.Headers.Get(key))
	}

	return nil
}

type BucketGetter interface {
	GetBucket() string
}

func joinURIPath(a, b string) string {
	if len(a) == 0 {
		a = "/"
	} else if a[0] != '/' {
		a = "/" + a
	}

	if len(b) != 0 && b[0] == '/' {
		b = b[1:]
	}

	if len(b) != 0 && a[len(a)-1] != '/' {
		a += "/"
	}

	return a + b
}
