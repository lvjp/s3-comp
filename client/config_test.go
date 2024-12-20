package client

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultEndpointResolver_ResolveEndpoint(t *testing.T) {
	region := "fr-me"
	bucket := "my-bucket"
	awsEndpoint := "https://s3.amazonaws.com"
	myEndpoint := "http://my-endpoint.com"

	testCases := []struct {
		Args     EndpointParameters
		Expected string
	}{
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: nil, UsePathStyle: false}, awsEndpoint},
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: nil, UsePathStyle: true}, awsEndpoint},
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: &myEndpoint, UsePathStyle: false}, myEndpoint},
		{EndpointParameters{Bucket: nil, Region: nil, Endpoint: &myEndpoint, UsePathStyle: true}, myEndpoint},
		{EndpointParameters{Bucket: nil, Region: &region, Endpoint: nil, UsePathStyle: false}, "https://s3.fr-me.amazonaws.com"},
		{EndpointParameters{Bucket: nil, Region: &region, Endpoint: nil, UsePathStyle: true}, "https://s3.fr-me.amazonaws.com"},
		{EndpointParameters{Bucket: nil, Region: &region, Endpoint: &myEndpoint, UsePathStyle: false}, myEndpoint},
		{EndpointParameters{Bucket: nil, Region: &region, Endpoint: &myEndpoint, UsePathStyle: true}, myEndpoint},
		{EndpointParameters{Bucket: &bucket, Region: nil, Endpoint: nil, UsePathStyle: false}, "https://my-bucket.s3.amazonaws.com"},
		{EndpointParameters{Bucket: &bucket, Region: nil, Endpoint: nil, UsePathStyle: true}, "https://s3.amazonaws.com/my-bucket"},
		{EndpointParameters{Bucket: &bucket, Region: nil, Endpoint: &myEndpoint, UsePathStyle: false}, "http://my-bucket.my-endpoint.com"},
		{EndpointParameters{Bucket: &bucket, Region: nil, Endpoint: &myEndpoint, UsePathStyle: true}, "http://my-endpoint.com/my-bucket"},
		{EndpointParameters{Bucket: &bucket, Region: &region, Endpoint: nil, UsePathStyle: false}, "https://my-bucket.s3.fr-me.amazonaws.com"},
		{EndpointParameters{Bucket: &bucket, Region: &region, Endpoint: nil, UsePathStyle: true}, "https://s3.fr-me.amazonaws.com/my-bucket"},
		{EndpointParameters{Bucket: &bucket, Region: &region, Endpoint: &myEndpoint, UsePathStyle: false}, "http://my-bucket.my-endpoint.com"},
		{EndpointParameters{Bucket: &bucket, Region: &region, Endpoint: &myEndpoint, UsePathStyle: true}, "http://my-endpoint.com/my-bucket"},
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
