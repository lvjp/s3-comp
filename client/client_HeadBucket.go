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

func (input *HeadBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *HeadBucketInput) MarshalHTTP(ctx context.Context, req *http.Request) error {
	req.Method = http.MethodHead

	if input.ExpectedBucketOwner != nil {
		req.Header.Set("X-Amz-Expected-Bucket-Owner", *input.ExpectedBucketOwner)
	}

	return nil
}

type HeadBucketOutput struct {
	BucketRegion     *string
	AccessPointAlias *string
}

func (output *HeadBucketOutput) UnmarshalHTTP(ctx context.Context, resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return NewAPIResponseError(ctx, resp)
	}

	extract := func(key string) *string {
		values := resp.Header.Values(key)
		if values == nil {
			return nil
		}

		return &values[0]
	}

	output.BucketRegion = extract("X-Amz-Bucket-Region")
	output.AccessPointAlias = extract("X-Amz-Access-Point-Alias")

	return nil
}

const operationHeadBucket = "HeadBucket"

func (c *Client) HeadBucket(ctx context.Context, input *HeadBucketInput) (*HeadBucketOutput, error) {
	ctx = withOperationName(ctx, operationHeadBucket)

	if input.Bucket == "" {
		return nil, NewSDKErrorBucketIsMandatory(ctx)
	}

	req := newRequest()

	if err := input.MarshalHTTP(ctx, req); err != nil {
		return nil, NewSDKError(ctx, err.Error())
	}

	if err := c.resolve(ctx, req, input); err != nil {
		return nil, NewSDKError(ctx, err.Error())
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, NewAPITransportError(ctx, err)
	}
	defer resp.Body.Close()

	output := new(HeadBucketOutput)
	if err := output.UnmarshalHTTP(ctx, resp); err != nil {
		return nil, err
	}

	return output, nil
}
