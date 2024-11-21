package client

import (
	"context"
	"net/http"
)

type HeadBucketInput struct {
	// bucket is required
	Bucket string

	ExpectedBucketOwner *string
}

type HeadBucketOutput struct {
	BucketRegion     *string
	AccessPointAlias *string
}

const operationHeadBucket = "HeadBucket"

func (c *Client) HeadBucket(ctx context.Context, input *HeadBucketInput) (*HeadBucketOutput, error) {
	ctx = withOperationName(ctx, operationHeadBucket)

	if input.Bucket == "" {
		return nil, NewSDKErrorBucketIsMandatory(ctx)
	}

	req, err := c.newRequest(ctx, &input.Bucket)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodHead

	if input.ExpectedBucketOwner != nil {
		req.Header.Set("X-Amz-Expected-Bucket-Owner", *input.ExpectedBucketOwner)
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, NewAPITransportError(ctx, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewAPIResponseError(ctx, resp)
	}

	extract := func(key string) *string {
		values := resp.Header.Values(key)
		if values == nil {
			return nil
		}

		return &values[0]
	}

	output := &HeadBucketOutput{
		BucketRegion:     extract("X-Amz-Bucket-Region"),
		AccessPointAlias: extract("X-Amz-Access-Point-Alias"),
	}

	return output, nil
}
