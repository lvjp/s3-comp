package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/lvjp/s3-comp/client/internal/pipeline"
)

type CreateBucketInput struct {
	// Bucket is required
	Bucket                     string
	ACL                        *string
	GrantFullControl           *string
	GrantRead                  *string
	GrantReadACP               *string
	GrantWrite                 *string
	GrantWriteACP              *string
	ObjectLockEnabledForBucket *string
	ObjectOwnership            *string

	CreateBucketConfiguration *CreateBucketConfiguration
}

func (input *CreateBucketInput) MarshalHTTP(ctx context.Context, req *http.Request) error {
	req.Method = http.MethodPut

	for key, val := range map[string]*string{
		"X-Amz-Acl":                        input.ACL,
		"X-Amz-Grant-Full-Control":         input.GrantFullControl,
		"X-Amz-Grant-Read":                 input.GrantRead,
		"X-Amz-Grant-Read-Acp":             input.GrantReadACP,
		"X-Amz-Grant-Write":                input.GrantWrite,
		"X-Amz-Grant-Write-Acp":            input.GrantWriteACP,
		"X-Amz-Bucket-Object-Lock-Enabled": input.ObjectLockEnabledForBucket,
		"X-Amz-Object-Ownership":           input.ObjectOwnership,
	} {
		if val != nil {
			req.Header.Set(key, *val)
		}
	}

	var hash [sha256.Size]byte
	if input.CreateBucketConfiguration != nil {
		inputBody, err := xml.Marshal(input.CreateBucketConfiguration)
		if err != nil {
			return err
		}

		req.ContentLength = int64(len(inputBody))
		req.Body = io.NopCloser(bytes.NewReader(inputBody))

		hash = sha256.Sum256(inputBody)
	} else {
		hash = sha256.Sum256(nil)
	}
	req.Header.Set("X-Amz-Content-Sha256", hex.EncodeToString(hash[:]))

	return nil
}

type CreateBucketOutput struct {
	Location *string
}

func (output *CreateBucketOutput) UnmarshalHTTP(ctx context.Context, resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return NewAPIResponseError(ctx, resp)
	}

	location := resp.Header.Values("Location")
	if location != nil {
		output.Location = &location[0]
	}

	return nil
}

func (c *Client) CreateBucket(ctx context.Context, input *CreateBucketInput) (*CreateBucketOutput, error) {
	const operationCreateBucket = "CreateBucket"

	output := new(CreateBucketOutput)

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
		Context:  withOperationName(ctx, operationCreateBucket),
		Bucket:   &input.Bucket,
		S3Input:  input,
		S3Output: output,
	}

	if err := handler.Handle(mwCtx); err != nil {
		return nil, err
	}

	return output, nil
}
