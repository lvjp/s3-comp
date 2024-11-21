package client

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_CreateBucket(t *testing.T) {
	tc := ActionTestRunner[CreateBucketInput, CreateBucketOutput]{
		OperationName: "CreateBucket",
		MissingBucket: func() *CreateBucketInput {
			return &CreateBucketInput{}
		},
		Normal: func() (*CreateBucketInput, *CreateBucketOutput, http.HandlerFunc) {
			return &CreateBucketInput{
					Bucket: "TheBucket",
				},
				&CreateBucketOutput{
					Location: ToPointer("TheLocation"),
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Location", "TheLocation")
				}
		},
		Executor: func(c *Client) func(context.Context, *CreateBucketInput) (*CreateBucketOutput, error) {
			return c.CreateBucket
		},
	}

	tc.Run(t)
}
