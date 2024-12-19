package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const defaultRegion = "eu-west-1"
const defaultEndpoint = "https://s3.eu-west-1.amazonaws.com"

type Config struct {
	HTTPClient HTTPClient
	UserAgent  *string

	Region           string
	Endpoint         string
	EndpointResolver EndpointResolver
	UsePathStyle     bool
}

func (cfg *Config) SetDefaults() error {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	if (cfg.Region != "" && cfg.Endpoint == "") || (cfg.Region == "" && cfg.Endpoint != "") {
		return errors.New("region and endpoint should be both either nil or not nil")
	}

	if cfg.Region == "" {
		cfg.Region = defaultRegion
		cfg.Endpoint = defaultEndpoint
	}

	if cfg.EndpointResolver == nil {
		cfg.EndpointResolver = &DefaultEndpointResolver{}
	}

	return nil
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type EndpointParameters struct {
	Bucket       *string
	Region       *string
	Endpoint     *string
	UsePathStyle bool
}

type Endpoint struct {
	URI     url.URL
	Headers http.Header
}

type EndpointResolver interface {
	ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error)
}

type DefaultEndpointResolver struct {
}

func (der *DefaultEndpointResolver) ResolveEndpoint(_ context.Context, params EndpointParameters) (*Endpoint, error) {
	base, err := der.baseURI(&params)
	if err != nil {
		return nil, err
	}

	if params.Bucket != nil {
		if params.UsePathStyle {
			base.Path = joinURIPath(base.Path, *params.Bucket)
			base.RawPath = joinURIPath(base.RawPath, *params.Bucket)
		} else {
			base.Host = *params.Bucket + "." + base.Host
		}
	}

	return &Endpoint{
			URI:     *base,
			Headers: make(http.Header),
		},
		nil
}

func (*DefaultEndpointResolver) baseURI(params *EndpointParameters) (*url.URL, error) {
	var endpoint string
	if params.Endpoint == nil {
		if params.Region == nil {
			endpoint = "https://s3.amazonaws.com"
		} else {
			endpoint = fmt.Sprintf("https://s3.%s.amazonaws.com", *params.Region)
		}
	} else {
		endpoint = *params.Endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf(`endpoint scheme should be either "http" or "https", not %q`, u.Scheme)
	}

	if u.Host == "" {
		return nil, errors.New("endpoint should have a host configured")
	}

	return u, nil
}
