package client

import (
	"context"
	"net/http"
	"strconv"
)

var _ error = &SDKError{}

// SDKError is returned when an error occurs within the SDK.
type SDKError struct {
	OperationName string
	Message       string
}

func (err *SDKError) Error() string {
	return err.OperationName + ": " + err.Message
}

func NewSDKError(ctx context.Context, message string) *SDKError {
	return &SDKError{
		OperationName: getOperationName(ctx),
		Message:       message,
	}
}

func NewSDKErrorBucketIsMandatory(ctx context.Context) *SDKError {
	return NewSDKError(ctx, "bucket is mandatory")
}

// ServerError is returned when the HTTP call failed or the server emit an
// error response.
type ServerError struct {
	OperationName string
	message       *string
	response      *http.Response
}

func (err *ServerError) Error() string {
	if err.message != nil {
		return err.OperationName + ": " + *err.message
	}

	return err.OperationName + ": unexpected http line " + strconv.Quote(err.response.Status)
}

func NewAPITransportError(ctx context.Context, message string) *ServerError {
	return &ServerError{
		OperationName: getOperationName(ctx),
		message:       &message,
	}
}

func NewAPIResponseError(ctx context.Context, resp *http.Response) *ServerError {
	return &ServerError{
		OperationName: getOperationName(ctx),
		response:      resp,
	}
}
