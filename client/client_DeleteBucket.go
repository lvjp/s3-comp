package client

import (
	"context"
	"net/http"

	"github.com/lvjp/s3-comp/client/internal/pipeline"
)

type DeleteBucketInput struct {
	// bucket is required
	Bucket string

	ExpectedBucketOwner *string
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

func (c *Client) DeleteBucket(ctx context.Context, input *DeleteBucketInput) (*DeleteBucketOutput, error) {
	const operationDeleteBucket = "DeleteBucket"

	output := new(DeleteBucketOutput)

	handler := pipeline.NewPipeline(
		pipeline.HandlerFunc(c.doRequest),
		mandatoryBucketMiddleware,
		initHTTPRequestMiddleware,
		httpMarshalerMiddleware,
		c.resolveMiddleware,
		httpUnmarshalerMiddleware,
	)

	mwCtx := &pipeline.MiddlewareContext{
		Context:  withOperationName(ctx, operationDeleteBucket),
		Bucket:   &input.Bucket,
		S3Input:  input,
		S3Output: output,
	}

	if err := handler.Handle(mwCtx); err != nil {
		return nil, err
	}

	return output, nil
}
