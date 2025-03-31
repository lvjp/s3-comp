package client

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/lvjp/s3-comp/client/types"

	"github.com/valyala/fasthttp"
)

type ListBucketsInput struct {
	BucketRegion      *string
	ContinuationToken *string
	MaxBuckets        *string
	Prefix            *string
}

func (input *ListBucketsInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodGet)

	args := req.URI().QueryArgs()
	setQuery(args, QueryBucketRegion, input.BucketRegion)
	setQuery(args, QueryContinuationToken, input.ContinuationToken)
	setQuery(args, QueryMaxBuckets, input.MaxBuckets)
	setQuery(args, QueryPrefix, input.Prefix)

	return nil
}

type ListBucketsOutput struct {
	Payload *types.ListAllMyBucketsResult
}

func (output *ListBucketsOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if resp.StatusCode() != fasthttp.StatusOK {
		return NewServerSideError(resp)
	}

	var payload types.ListAllMyBucketsResult
	if err := xml.Unmarshal(resp.Body(), &payload); err != nil {
		return fmt.Errorf("ListBuckets: cannot parse response body: %w", err)
	}
	output.Payload = &payload
	return nil
}

func (c *Client) ListBuckets(ctx context.Context, input *ListBucketsInput, optFns ...func(*Options)) (*ListBucketsOutput, *Metadata, error) {
	return PerformCall[*ListBucketsInput, *ListBucketsOutput](ctx, c, input, optFns...)
}
