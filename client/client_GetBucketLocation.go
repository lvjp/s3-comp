package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"net/http"

	"github.com/lvjp/s3-comp/client/internal/pipeline"
)

type GetBucketLocationInput struct {
	// bucket is required
	Bucket         string
	ExpectedBucket *string
}

func (input *GetBucketLocationInput) MarshalHTTP(_ context.Context, req *http.Request) error {
	req.Method = http.MethodGet

	q := req.URL.Query()
	q.Set("location", "")
	req.URL.RawQuery = q.Encode()

	if input.ExpectedBucket != nil {
		req.Header.Set("X-Amz-Expected-Bucket-Owner", *input.ExpectedBucket)
	}

	hash := sha256.Sum256(nil)
	req.Header.Set("X-Amz-Content-Sha256", hex.EncodeToString(hash[:]))

	return nil
}

type GetBucketLocationOutput struct {
	XMLName            xml.Name            `xml:"http://s3.amazonaws.com/doc/2006-03-01/ LocationConstraint"`
	LocationConstraint *LocationConstraint `xml:",chardata"`
}

func (output *GetBucketLocationOutput) UnmarshalHTTP(ctx context.Context, resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return NewAPIResponseError(ctx, resp)
	}

	return xml.NewDecoder(resp.Body).Decode(output)
}

func (c *Client) GetBucketLocation(ctx context.Context, input *GetBucketLocationInput) (*GetBucketLocationOutput, error) {
	const operationName = "GetBucketLocation"

	output := new(GetBucketLocationOutput)

	handler := pipeline.NewPipeline(
		pipeline.HandlerFunc(c.doRequest),
		mandatoryBucketMiddleware,
		initHTTPRequestMiddleware,
		httpMarshalerMiddleware,
		userAgentMiddleware(c.config.UserAgent),
		disableDefaultTransportGzipCompression,
		c.resolveMiddleware,
		httpUnmarshalerMiddleware,
	)

	mwCtx := &pipeline.MiddlewareContext{
		Context:  withOperationName(ctx, operationName),
		Bucket:   &input.Bucket,
		S3Input:  input,
		S3Output: output,
	}

	if err := handler.Handle(mwCtx); err != nil {
		return nil, err
	}

	return output, nil
}
