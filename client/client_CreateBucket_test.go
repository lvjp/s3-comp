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
			location := "TheLocation"
			return &CreateBucketInput{
					Bucket: "TheBucket",
				},
				&CreateBucketOutput{
					Location: &location,
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Location", location)
				}
		},
		Executor: func(c *Client) func(context.Context, *CreateBucketInput) (*CreateBucketOutput, error) {
			return c.CreateBucket
		},
	}

	tc.Run(t)
}
