package client

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_HeadBucket(t *testing.T) {
	tc := ActionTestRunner[HeadBucketInput, HeadBucketOutput]{
		OperationName: "HeadBucket",
		MissingBucket: func() *HeadBucketInput {
			return &HeadBucketInput{}
		},
		Normal: func() (*HeadBucketInput, *HeadBucketOutput, http.HandlerFunc) {
			return &HeadBucketInput{
					Bucket: "TheBucket",
				},
				&HeadBucketOutput{
					BucketRegion:     ToPointer("BucketRegion"),
					AccessPointAlias: ToPointer("AccessPointAlias"),
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Amz-Bucket-Region", "BucketRegion")
					w.Header().Set("X-Amz-Access-Point-Alias", "AccessPointAlias")
				}
		},
		Executor: func(c *Client) func(context.Context, *HeadBucketInput) (*HeadBucketOutput, error) {
			return c.HeadBucket
		},
	}

	tc.Run(t)
}
