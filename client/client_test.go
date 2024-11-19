package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type testResolver struct{}

func (*testResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error) {
	resolver := &DefaultEndpointResolver{}

	endpoint, err := resolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return nil, err
	}

	decoded, err := url.Parse(*params.Endpoint)
	if err != nil {
		return nil, err
	}

	// Little trick for DNS
	endpoint.Headers.Set("Host", endpoint.URI.Host)
	endpoint.URI.Host = decoded.Host

	return endpoint, nil
}

func TestNew(t *testing.T) {
	client, err := New(Config{})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Same(t, http.DefaultClient, client.config.HTTPClient)
}
