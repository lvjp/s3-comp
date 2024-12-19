package client

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"net/http"
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

func (input *CreateBucketInput) GetBucket() string {
	return input.Bucket
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

	if input.CreateBucketConfiguration != nil {
		inputBody, err := xml.Marshal(input.CreateBucketConfiguration)
		if err != nil {
			return err
		}

		req.ContentLength = int64(len(inputBody))
		req.Body = io.NopCloser(bytes.NewReader(inputBody))
	}

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

const operationCreateBucket = "CreateBucket"

func (c *Client) CreateBucket(ctx context.Context, input *CreateBucketInput) (*CreateBucketOutput, error) {
	ctx = withOperationName(ctx, operationCreateBucket)

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

	output := new(CreateBucketOutput)
	if err := output.UnmarshalHTTP(ctx, resp); err != nil {
		return nil, err
	}

	return output, nil
}
