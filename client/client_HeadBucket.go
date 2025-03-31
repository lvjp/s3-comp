package client

import (
	"context"

	"github.com/valyala/fasthttp"
)

var _ RequiredBucketInterface = (*HeadBucketInput)(nil)

type HeadBucketInput struct {
	// Bucket is mandatory
	Bucket string

	ExpectedBucketOwner *string
}

func (input *HeadBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *HeadBucketInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodHead)

	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)

	return nil
}

type HeadBucketOutput struct {
	AccessPointAlias *string
	BucketRegion     *string
}

func (output *HeadBucketOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if resp.StatusCode() != fasthttp.StatusOK {
		return NewServerSideError(resp)
	}

	output.BucketRegion = extractHeader(&resp.Header, HeaderXAmzBucketRegion)
	output.AccessPointAlias = extractHeader(&resp.Header, HeaderXAmzAccessPointAlias)

	return nil
}

func (c *Client) HeadBucket(ctx context.Context, input *HeadBucketInput, optFns ...func(*Options)) (*HeadBucketOutput, *Metadata, error) {
	return PerformCall[*HeadBucketInput, *HeadBucketOutput](ctx, c, input, optFns...)
}
