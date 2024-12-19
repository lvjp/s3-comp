package client

import (
	"context"
	"net/http"
)

type DeleteBucketInput struct {
	// bucket is required
	Bucket string

	ExpectedBucketOwner *string
}

func (input *DeleteBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *DeleteBucketInput) MarshalHTTP(ctx context.Context, req *http.Request) error {
	req.Method = http.MethodDelete
	req.Body = http.NoBody

	if input.ExpectedBucketOwner != nil {
		req.Header.Set("X-Amz-Expected-Bucket-Owner", *input.ExpectedBucketOwner)
	}

	return nil
}

type DeleteBucketOutput struct {
}

func (*DeleteBucketOutput) UnmarshalHTTP(ctx context.Context, resp *http.Response) error {
	if resp.StatusCode != http.StatusNoContent {
		return NewAPIResponseError(ctx, resp)
	}

	return nil
}

const operationDeleteBucket = "DeleteBucket"

func (c *Client) DeleteBucket(ctx context.Context, input *DeleteBucketInput) (*DeleteBucketOutput, error) {
	ctx = withOperationName(ctx, operationDeleteBucket)

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

	output := new(DeleteBucketOutput)
	if err := output.UnmarshalHTTP(ctx, resp); err != nil {
		return nil, err
	}

	return &DeleteBucketOutput{}, nil
}
