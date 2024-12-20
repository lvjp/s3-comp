package client

import (
	"context"
	"net/http"

	"github.com/lvjp/s3-comp/client/internal/pipeline"
)

type HeadBucketInput struct {
	// bucket is required
	Bucket string

	ExpectedBucketOwner *string
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

func (c *Client) HeadBucket(ctx context.Context, input *HeadBucketInput) (*HeadBucketOutput, error) {
	const operationHeadBucket = "HeadBucket"

	output := new(HeadBucketOutput)

	handler := pipeline.NewPipeline(
		pipeline.HandlerFunc(c.doRequest),
		mandatoryBucketMiddleware,
		initHTTPRequestMiddleware,
		httpMarshalerMiddleware,
		c.resolveMiddleware,
		httpUnmarshalerMiddleware,
	)

	mwCtx := &pipeline.MiddlewareContext{
		Context:  withOperationName(ctx, operationHeadBucket),
		Bucket:   &input.Bucket,
		S3Input:  input,
		S3Output: output,
	}

	if err := handler.Handle(mwCtx); err != nil {
		return nil, err
	}

	return output, nil
}
