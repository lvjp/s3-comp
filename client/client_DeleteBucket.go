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

type DeleteBucketOutput struct {
}

const operationDeleteBucket = "DeleteBucket"

func (c *Client) DeleteBucket(ctx context.Context, input *DeleteBucketInput) (*DeleteBucketOutput, error) {
	ctx = withOperationName(ctx, operationDeleteBucket)

	if input.Bucket == "" {
		return nil, NewSDKErrorBucketIsMandatory(ctx)
	}

	req, err := c.newRequest(ctx, &input.Bucket)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodDelete

	if input.ExpectedBucketOwner != nil {
		req.Header.Set("X-Amz-Expected-Bucket-Owner", *input.ExpectedBucketOwner)
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, NewAPITransportError(ctx, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return nil, NewAPIResponseError(ctx, resp)
	}

	return &DeleteBucketOutput{}, nil
}
