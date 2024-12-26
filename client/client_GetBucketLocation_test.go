package client

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetBucketLocation(t *testing.T) {
	tc := ActionTestRunner[GetBucketLocationInput, GetBucketLocationOutput, s3.GetBucketLocationInput]{
		OperationName: "GetBucketLocation",
		MissingBucket: func() *GetBucketLocationInput {
			return &GetBucketLocationInput{}
		},
		Normal: func() (*GetBucketLocationInput, *GetBucketLocationOutput, *s3.GetBucketLocationInput, func(*testing.T) http.HandlerFunc) {
			bucket := "the-bucket"
			expectedOwner := "expected-owner"
			bucketLocation := LocationConstraint("bucket-location")

			return &GetBucketLocationInput{
					Bucket:         bucket,
					ExpectedBucket: &expectedOwner,
				},
				&GetBucketLocationOutput{
					XMLName: xml.Name{
						Space: "http://s3.amazonaws.com/doc/2006-03-01/",
						Local: "LocationConstraint",
					},
					LocationConstraint: &bucketLocation,
				},
				&s3.GetBucketLocationInput{
					ExpectedBucketOwner: &expectedOwner,
					Bucket:              &bucket,
				},
				func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						body := []byte(xml.Header + `<LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">` + bucketLocation + "</LocationConstraint>")
						w.Header().Set("Content-Length", strconv.Itoa(len(body)))
						w.Header().Set("Content-Type", "text/xml")
						w.WriteHeader(http.StatusOK)
						n, err := w.Write(body)
						if assert.NoError(t, err) {
							assert.Equal(t, len(body), n)
						}
					}
				}
		},
		Executor: func(c *Client) func(context.Context, *GetBucketLocationInput) (*GetBucketLocationOutput, error) {
			return c.GetBucketLocation
		},
		AWSExecute: func(c *s3.Client, ctx context.Context, input *s3.GetBucketLocationInput) error {
			_, err := c.GetBucketLocation(ctx, input)
			return err
		},
	}

	tc.Run(t)
}
