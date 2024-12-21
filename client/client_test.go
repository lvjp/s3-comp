package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c, err := New(Config{})
	require.NoError(t, err)
	require.NotNil(t, c)
}

type TestRunner interface {
	Name() string
	Run(t *testing.T)
}

type ActionTestRunner[Input, Output any] struct {
	OperationName string
	MissingBucket func() *Input
	MissingKey    func() *Input
	Normal        func() (*Input, *Output, func(*testing.T) http.HandlerFunc)
	Executor      func(*Client) func(context.Context, *Input) (*Output, error)
}

func (atr *ActionTestRunner[Input, Output]) Name() string {
	return atr.OperationName
}

func (atr *ActionTestRunner[Input, Output]) Run(t *testing.T) {
	atr.runValidation(t)
	atr.runNormal(t)
	atr.runErrors(t)
}

func (atr *ActionTestRunner[Input, Output]) runNormal(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		input, expected, handler := atr.Normal()

		ts := httptest.NewServer(handler(t))
		defer ts.Close()

		config := Config{
			HTTPClient:       ts.Client(),
			Region:           "fr-dev",
			Endpoint:         ts.URL,
			EndpointResolver: &hostHeaderResolver{},
		}

		c, err := New(config)
		require.NoError(t, err)

		actual, err := atr.Executor(c)(context.Background(), input)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func (atr *ActionTestRunner[Input, Output]) runErrors(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("server", func(t *testing.T) {
			httpClient := &error500HTTPClient{}
			httpResponse := httpClient.NewResponse()
			defer httpResponse.Body.Close()

			input, _, _ := atr.Normal()
			expectedErr := NewAPIResponseError(
				withOperationName(context.Background(), atr.OperationName),
				httpResponse,
			)

			c, err := New(Config{HTTPClient: httpClient})
			require.NoError(t, err)

			actualOutput, actualErr := atr.Executor(c)(context.Background(), input)
			require.EqualError(t, actualErr, expectedErr.Error())
			require.Nil(t, actualOutput)
		})

		t.Run("transport", func(t *testing.T) {
			httpClient := &transportErrorHTTPClient{}

			input, _, _ := atr.Normal()
			expectedErr := NewAPITransportError(
				withOperationName(context.Background(), atr.OperationName),
				httpClient.NewError(),
			)

			c, err := New(Config{
				HTTPClient: httpClient,
			})
			require.NoError(t, err)

			actualOutput, actualErr := atr.Executor(c)(context.Background(), input)
			require.EqualError(t, actualErr, expectedErr.Error())
			require.Nil(t, actualOutput)
		})
	})
}

func (atr *ActionTestRunner[Input, Output]) runValidation(t *testing.T) {
	t.Run("validation", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    func() *Input
			expected string
		}{
			{name: "missing_bucket", input: atr.MissingBucket, expected: ": bucket is mandatory"},
			{name: "missing_key", input: atr.MissingKey, expected: ": key is mandatory"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.input == nil {
					t.SkipNow()
				}

				c, err := New(Config{
					HTTPClient: &transportErrorHTTPClient{},
				})
				require.NoError(t, err)

				actual, err := atr.Executor(c)(context.Background(), tc.input())
				require.ErrorContains(t, err, tc.expected)
				require.Nil(t, actual)
			})
		}
	})
}

type transportErrorHTTPClient struct{}

func (*transportErrorHTTPClient) NewError() error {
	return errors.New("transportErrorHTTPClient")
}

func (t *transportErrorHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, t.NewError()
}

type error500HTTPClient struct{}

func (*error500HTTPClient) NewResponse() *http.Response {
	const statusCode = http.StatusInternalServerError
	return &http.Response{
		Status: fmt.Sprintf(
			"%d %s",
			statusCode,
			http.StatusText(statusCode),
		),
		StatusCode: statusCode,
		Body:       http.NoBody,
	}
}

func (e *error500HTTPClient) Do(*http.Request) (*http.Response, error) {
	return e.NewResponse(), nil
}

type hostHeaderResolver struct{}

func (*hostHeaderResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error) {
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
