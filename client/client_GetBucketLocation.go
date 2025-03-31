package client

import (
	"context"
	"encoding/xml"

	"github.com/lvjp/s3-comp/client/types"

	"github.com/valyala/fasthttp"
)

type GetBucketLocationInput struct {
	// Bucket is mandatory
	Bucket string

	ExpectedBucketOwner *string
}

func (input *GetBucketLocationInput) GetBucket() string {
	return input.Bucket
}

func (input *GetBucketLocationInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodGet)

	req.URI().QueryArgs().SetNoValue(QueryLocation)

	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)

	return nil
}

type GetBucketLocationOutput struct {
	XMLName            xml.Name                  `xml:"LocationConstraint"`
	LocationConstraint *types.LocationConstraint `xml:",chardata"`
}

func (output *GetBucketLocationOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if resp.StatusCode() != fasthttp.StatusOK {
		return NewServerSideError(resp)
	}

	return xml.Unmarshal(resp.Body(), output)
}

func (c *Client) GetBucketLocation(ctx context.Context, input *GetBucketLocationInput, optFns ...func(*Options)) (*GetBucketLocationOutput, *Metadata, error) {
	return PerformCall[*GetBucketLocationInput, *GetBucketLocationOutput](ctx, c, input, optFns...)
}
