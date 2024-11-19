package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_CreateBucket(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		t.Run("bucket", func(t *testing.T) {
			c, err := New(Config{})
			require.NoError(t, err)
			_, err = c.CreateBucket(context.Background(), &CreateBucketInput{})
			require.EqualError(t, err, "CreateBucket: bucket is mandatory")
		})
	})

	t.Run("request", func(t *testing.T) {
		t.Run("500", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			conf := Config{
				Region:           "fr-dev",
				Endpoint:         ts.URL,
				EndpointResolver: &testResolver{},
			}

			c, err := New(conf)
			require.NoError(t, err)

			createBucketOutput, err := c.CreateBucket(
				context.Background(),
				&CreateBucketInput{
					Bucket: "my-pouet-lv123",
				},
			)
			require.EqualError(t, err, `CreateBucket: unexpected http line "500 Internal Server Error"`)
			require.Nil(t, createBucketOutput)
		})

		t.Run("VHostStyle", func(t *testing.T) {
			var requestDump []byte

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var err error
				requestDump, err = httputil.DumpRequest(r, true)
				if err != nil {
					t.Error("Request dump error:", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Location", "TheLocation")
			}))
			defer ts.Close()

			conf := Config{
				HTTPClient: &http.Client{
					Transport: &http.Transport{
						DisableCompression: true,
					},
				},
				Region:           "fr-dev",
				Endpoint:         ts.URL,
				EndpointResolver: &testResolver{},
			}

			c, err := New(conf)
			require.NoError(t, err)

			createBucketOutput, err := c.CreateBucket(
				context.Background(),
				&CreateBucketInput{
					Bucket: "my-pouet-lv123",
				},
			)
			require.NoError(t, err)

			expectedRequest := "PUT / HTTP/1.1\r\n" +
				"Host: " + ts.Listener.Addr().String() + "\r\n" +
				"Content-Length: 0\r\n" +
				"User-Agent: Go-http-client/1.1\r\n" +
				"\r\n"
			require.Equal(t, []byte(expectedRequest), requestDump)

			expectedOutput := &CreateBucketOutput{
				Location: ToPointer("TheLocation"),
			}
			require.Equal(t, expectedOutput, createBucketOutput)
		})

		t.Run("PathStyle", func(t *testing.T) {
			var requestDump []byte

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var err error
				requestDump, err = httputil.DumpRequest(r, true)
				if err != nil {
					t.Error("Request dump error:", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Location", "TheLocation")
			}))
			defer ts.Close()

			conf := Config{
				HTTPClient: &http.Client{
					Transport: &http.Transport{
						DisableCompression: true,
					},
				},
				Region:           "fr-dev",
				Endpoint:         ts.URL,
				EndpointResolver: &testResolver{},
				UsePathStyle:     true,
			}

			c, err := New(conf)
			require.NoError(t, err)

			createBucketOutput, err := c.CreateBucket(
				context.Background(),
				&CreateBucketInput{
					Bucket: "my-pouet-lv123",
				},
			)
			require.NoError(t, err)

			expectedRequest := "PUT /my-pouet-lv123 HTTP/1.1\r\n" +
				"Host: " + ts.Listener.Addr().String() + "\r\n" +
				"Content-Length: 0\r\n" +
				"User-Agent: Go-http-client/1.1\r\n" +
				"\r\n"
			require.Equal(t, []byte(expectedRequest), requestDump)

			expectedOutput := &CreateBucketOutput{
				Location: ToPointer("TheLocation"),
			}
			require.Equal(t, expectedOutput, createBucketOutput)
		})
	})
}
