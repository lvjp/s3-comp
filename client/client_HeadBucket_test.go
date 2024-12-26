package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestClient_HeadBucket(t *testing.T) {
	tc := ActionTestRunner[HeadBucketInput, HeadBucketOutput, s3.HeadBucketInput]{
		OperationName: "HeadBucket",
		MissingBucket: func() *HeadBucketInput {
			return &HeadBucketInput{}
		},
		Normal: func() (*HeadBucketInput, *HeadBucketOutput, *s3.HeadBucketInput, func(t *testing.T) http.HandlerFunc) {
			bucket := "TheBucket"
			bucketRegion := "BucketRegion"
			accessPointAlias := "false"
			return &HeadBucketInput{
					Bucket: bucket,
				},
				&HeadBucketOutput{
					BucketRegion:     &bucketRegion,
					AccessPointAlias: &accessPointAlias,
				},
				&s3.HeadBucketInput{
					Bucket: &bucket,
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
		AWSExecute: func(c *s3.Client, ctx context.Context, input *s3.HeadBucketInput) error {
			_, err := c.HeadBucket(ctx, input)
			return err
		},
	}

	tc.Run(t)
}
