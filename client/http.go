package client

import "github.com/valyala/fasthttp"

type HTTPClient interface {
	Do(*fasthttp.Request, *fasthttp.Response) error
}

var _ HTTPClient = DefaultHTTPClient

var DefaultHTTPClient = &fasthttp.Client{
	NoDefaultUserAgentHeader: true,
	RetryIfErr: func(request *fasthttp.Request, attempts int, err error) (resetTimeout bool, retry bool) {
		return false, false
	},
}
