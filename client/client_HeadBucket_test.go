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
		Normal: func() (*HeadBucketInput, *HeadBucketOutput, func(t *testing.T) http.HandlerFunc) {
			bucketRegion := "BucketRegion"
			accessPointAlias := "AccessPointAlias"
			return &HeadBucketInput{
					Bucket: "TheBucket",
				},
				&HeadBucketOutput{
					BucketRegion:     &bucketRegion,
					AccessPointAlias: &accessPointAlias,
				},
				func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Amz-Bucket-Region", bucketRegion)
						w.Header().Set("X-Amz-Access-Point-Alias", accessPointAlias)
					}
				}
		},
		Executor: func(c *Client) func(context.Context, *HeadBucketInput) (*HeadBucketOutput, error) {
			return c.HeadBucket
		},
	}

	tc.Run(t)
}
