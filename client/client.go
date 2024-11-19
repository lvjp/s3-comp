package client

import (
	"context"
	"net/http"
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

func (c *Client) newRequest(ctx context.Context, bucket *string) (*http.Request, error) {
	endpoint, err := c.config.EndpointResolver.ResolveEndpoint(
		ctx,
		EndpointParameters{
			Bucket:       bucket,
			Region:       &c.config.Region,
			Endpoint:     &c.config.Endpoint,
			UsePathStyle: c.config.UsePathStyle,
		},
	)
	if err != nil {
		if _, ok := c.config.EndpointResolver.(*DefaultEndpointResolver); ok {
			return nil, err
		}

		return nil, NewSDKError(ctx, "cannot resolve endpoint: "+err.Error())
	}

	req := &http.Request{
		URL:    &endpoint.URI,
		Header: endpoint.Headers,
	}

	if c.config.UserAgent != nil {
		req.Header.Set("User-Agent", *c.config.UserAgent)
	}

	return req, nil
}
