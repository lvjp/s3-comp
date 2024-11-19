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

type CreateBucketOutput struct {
	Location *string
}

const operationCreateBucket = "CreateBucket"

func (c *Client) CreateBucket(ctx context.Context, input *CreateBucketInput) (*CreateBucketOutput, error) {
	ctx = withOperationName(ctx, operationCreateBucket)

	if input.Bucket == "" {
		return nil, NewSDKErrorBucketIsMandatory(ctx)
	}

	req, err := c.newRequest(ctx, &input.Bucket)
	if err != nil {
		return nil, err
	}

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
		inputBody, marshalErr := xml.Marshal(input.CreateBucketConfiguration)
		if marshalErr != nil {
			return nil, NewSDKError(ctx, marshalErr.Error())
		}

		req.ContentLength = int64(len(inputBody))
		req.Body = io.NopCloser(bytes.NewReader(inputBody))
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, NewAPITransportError(ctx, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewAPIResponseError(ctx, resp)
	}

	extract := func(key string) *string {
		values := resp.Header.Values(key)
		if values == nil {
			return nil
		}

		return &values[0]
	}

	output := &CreateBucketOutput{
		Location: extract("Location"),
	}

	return output, nil
}
