package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestClient_DeleteBucket(t *testing.T) {
	tc := ActionTestRunner[DeleteBucketInput, DeleteBucketOutput, s3.DeleteBucketInput]{
		OperationName: "DeleteBucket",
		MissingBucket: func() *DeleteBucketInput {
			return &DeleteBucketInput{}
		},
		Normal: func() (*DeleteBucketInput, *DeleteBucketOutput, *s3.DeleteBucketInput, func(t *testing.T) http.HandlerFunc) {
			bucket := "TheBucket"
			return &DeleteBucketInput{
					Bucket: bucket,
				},
				&DeleteBucketOutput{},
				&s3.DeleteBucketInput{
					Bucket: &bucket,
				},
				func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusNoContent)
					}
				}
		},
		Executor: func(c *Client) func(context.Context, *DeleteBucketInput) (*DeleteBucketOutput, error) {
			return c.DeleteBucket
		},
		AWSExecute: func(c *s3.Client, ctx context.Context, input *s3.DeleteBucketInput) error {
			_, err := c.DeleteBucket(ctx, input)
			return err
		},
	}

	tc.Run(t)
}
