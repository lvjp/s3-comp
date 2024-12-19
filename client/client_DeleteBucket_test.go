package client

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_DeleteBucket(t *testing.T) {
	tc := ActionTestRunner[DeleteBucketInput, DeleteBucketOutput]{
		OperationName: "DeleteBucket",
		MissingBucket: func() *DeleteBucketInput {
			return &DeleteBucketInput{}
		},
		Normal: func() (*DeleteBucketInput, *DeleteBucketOutput, http.HandlerFunc) {
			return &DeleteBucketInput{
					Bucket: "TheBucket",
				},
				&DeleteBucketOutput{},
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}
		},
		Executor: func(c *Client) func(context.Context, *DeleteBucketInput) (*DeleteBucketOutput, error) {
			return c.DeleteBucket
		},
	}

	tc.Run(t)
}
