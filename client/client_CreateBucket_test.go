package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestClient_CreateBucket(t *testing.T) {
	tc := ActionTestRunner[CreateBucketInput, CreateBucketOutput, s3.CreateBucketInput]{
		OperationName: "CreateBucket",
		MissingBucket: func() *CreateBucketInput {
			return &CreateBucketInput{}
		},
		Normal: func() (*CreateBucketInput, *CreateBucketOutput, *s3.CreateBucketInput, func(t *testing.T) http.HandlerFunc) {
			location := "TheLocation"
			bucket := "TheBucket"
			return &CreateBucketInput{
					Bucket: bucket,
				},
				&CreateBucketOutput{
					Location: &location,
				},
				&s3.CreateBucketInput{
					Bucket: &bucket,
				},
				func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Location", location)
					}
				}
		},
		Executor: func(c *Client) func(context.Context, *CreateBucketInput) (*CreateBucketOutput, error) {
			return c.CreateBucket
		},
		AWSExecute: func(c *s3.Client, ctx context.Context, input *s3.CreateBucketInput) error {
			_, err := c.CreateBucket(ctx, input)
			return err
		},
	}

	tc.Run(t)
}
