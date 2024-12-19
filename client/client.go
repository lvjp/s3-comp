package client

import (
	"bytes"
	"io"
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

func (c *Client) doRequest(ctx *MiddlewareContext) error {
	resp, err := c.config.HTTPClient.Do(ctx.HTTPRequest)
	if err != nil {
		return NewAPITransportError(ctx, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))
	ctx.HTTPResponse = resp

	return nil
}

func (c *Client) resolveMiddleware(next Handler) Handler {
	return HandlerFunc(func(ctx *MiddlewareContext) error {
		params := EndpointParameters{
			Region:       &c.config.Region,
			Endpoint:     &c.config.Endpoint,
			UsePathStyle: c.config.UsePathStyle,
			Bucket:       ctx.Bucket,
		}

		endpoint, err := c.config.EndpointResolver.ResolveEndpoint(ctx, params)
		if err != nil {
			if _, ok := c.config.EndpointResolver.(*DefaultEndpointResolver); ok {
				return err
			}

			return NewSDKError(ctx, "cannot resolve endpoint: "+err.Error())
		}

		ctx.HTTPRequest.URL.Scheme = endpoint.URI.Scheme
		ctx.HTTPRequest.URL.Host = endpoint.URI.Host
		ctx.HTTPRequest.URL.Path = joinURIPath(endpoint.URI.Path, ctx.HTTPRequest.URL.Path)
		ctx.HTTPRequest.URL.RawPath = joinURIPath(endpoint.URI.RawPath, ctx.HTTPRequest.URL.RawPath)

		for key := range endpoint.Headers {
			ctx.HTTPRequest.Header.Set(key, endpoint.Headers.Get(key))
		}

		return next.Handle(ctx)
	})
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
