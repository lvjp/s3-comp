package client

import (
	"context"

	"github.com/valyala/fasthttp"
)

type DeleteBucketInput struct {
	// bucket is required
	Bucket string

	ExpectedBucketOwner *string
}

func (input *DeleteBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *DeleteBucketInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodDelete)

	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)

	return nil
}

type DeleteBucketOutput struct {
}

func (*DeleteBucketOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if resp.StatusCode() != fasthttp.StatusNoContent {
		return NewServerSideError(resp)
	}

	return nil
}

func (c *Client) DeleteBucket(ctx context.Context, input *DeleteBucketInput, optFns ...func(*Options)) (*DeleteBucketOutput, *Metadata, error) {
	return PerformCall[*DeleteBucketInput, *DeleteBucketOutput](ctx, c, input, optFns...)
}
