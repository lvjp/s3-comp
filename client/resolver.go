package client

import (
	"context"
)

type EndpointParameters struct {
	Bucket       string
	Key          string
	Host         string
	UseSSL       bool
	UsePathStyle bool
}

type Endpoint struct {
	URL string
}

type EndpointResolver interface {
	ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error)
}

var DefaultEndpointResolver = &defaultEndpointResolver{}

type defaultEndpointResolver struct{}

func (*defaultEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error) {
	url := "http"
	if params.UseSSL {
		url += "s"
	}

	url += "://"

	if !params.UsePathStyle && params.Bucket != "" {
		url += params.Bucket
		url += "."
	}

	url += params.Host

	if params.UsePathStyle && params.Bucket != "" {
		url += "/"
		url += params.Bucket
	}

	if params.Key != "" {
		url += "/"
		url += params.Key
	}

	return &Endpoint{URL: url}, nil
}
