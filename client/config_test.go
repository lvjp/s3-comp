package client

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultEndpointResolver_ResolveEndpoint(t *testing.T) {
	testCases := []struct {
		Args     EndpointParameters
		Expected string
	}{
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: nil, UsePathStyle: false}, "https://s3.amazonaws.com"},
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: nil, UsePathStyle: true}, "https://s3.amazonaws.com"},
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: false}, "http://my-endpoint.com"},
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: true}, "http://my-endpoint.com"},
		{EndpointParameters{Bucket: nil, Region: ToPointer("fr-me"), Endpoint: nil, UsePathStyle: false}, "https://s3.fr-me.amazonaws.com"},
		{EndpointParameters{Bucket: nil, Region: ToPointer("fr-me"), Endpoint: nil, UsePathStyle: true}, "https://s3.fr-me.amazonaws.com"},
		{EndpointParameters{Bucket: nil, Region: ToPointer("fr-me"), Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: false}, "http://my-endpoint.com"},
		{EndpointParameters{Bucket: nil, Region: ToPointer("fr-me"), Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: true}, "http://my-endpoint.com"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: nil, Endpoint: nil, UsePathStyle: false}, "https://my-bucket.s3.amazonaws.com"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: nil, Endpoint: nil, UsePathStyle: true}, "https://s3.amazonaws.com/my-bucket"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: nil, Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: false}, "http://my-bucket.my-endpoint.com"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: nil, Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: true}, "http://my-endpoint.com/my-bucket"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: ToPointer("fr-me"), Endpoint: nil, UsePathStyle: false}, "https://my-bucket.s3.fr-me.amazonaws.com"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: ToPointer("fr-me"), Endpoint: nil, UsePathStyle: true}, "https://s3.fr-me.amazonaws.com/my-bucket"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: ToPointer("fr-me"), Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: false}, "http://my-bucket.my-endpoint.com"},
		{EndpointParameters{Bucket: ToPointer("my-bucket"), Region: ToPointer("fr-me"), Endpoint: ToPointer("http://my-endpoint.com"), UsePathStyle: true}, "http://my-endpoint.com/my-bucket"},
	}

	var resolver DefaultEndpointResolver

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := resolver.ResolveEndpoint(context.Background(), tc.Args)
			require.NoError(t, err)
			require.Equal(t, tc.Expected, actual.URI.String())
		})
	}
}
