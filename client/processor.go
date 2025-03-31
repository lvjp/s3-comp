package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	chain_of_responsibility "github.com/lvjp/s3-comp/client/pkg/design-patterns/chain-of-responsibility"

	"github.com/valyala/fasthttp"
)

type Handler[Input any, Output any] = chain_of_responsibility.Handler[*handlerInput[Input], *handlerOutput[Output]]

type HTTPRequestMarshaler interface {
	MarshalHTTP(req *fasthttp.Request) error
}

type HTTPResponseUnmarshaler interface {
	UnmarshalHTTP(resp *fasthttp.Response) error
}

type RequiredBucketInterface interface {
	GetBucket() string
}

var _ RequiredBucketInterface = (RequiredBucketKeyInterface)(nil)

type RequiredBucketKeyInterface interface {
	RequiredBucketInterface
	GetKey() string
}

type handlerInput[Input any] struct {
	Options       *Options
	CallInput     Input
	ServerRequest fasthttp.Request
}

type handlerOutput[Output any] struct {
	CallOutputV3   Output
	ServerResponse *fasthttp.Response
}

type httpRequesterHandler[Input any, Output any] struct{}

func (*httpRequesterHandler[Input, Output]) Handle(ctx context.Context, input *handlerInput[Input]) (*handlerOutput[Output], error) {
	output := &handlerOutput[Output]{
		ServerResponse: &fasthttp.Response{},
	}

	if err := input.Options.HTTPClient.Do(&input.ServerRequest, output.ServerResponse); err != nil {
		return nil, fmt.Errorf("HTTP request error: %v", err)
	}

	return output, nil
}

type errorMiddleware[Input any, Output any] struct{}

func (*errorMiddleware[Input, Output]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, Output]) (*handlerOutput[Output], error) {
	output, err := next.Handle(ctx, input)
	if err == nil {
		return output, nil
	}

	var serverSideError *ServerSideError
	var clientSideError *ClientSideError

	if !errors.As(err, &serverSideError) && !errors.As(err, &clientSideError) {
		err = &ClientSideError{Err: err}
	}

	return output, err
}

type configValidationMiddleware[Input any, Output any] struct{}

func (*configValidationMiddleware[Input, Output]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, Output]) (*handlerOutput[Output], error) {
	if err := input.Options.validate(); err != nil {
		return nil, err
	}

	return next.Handle(ctx, input)
}

type userAgentMiddleware[Input any, Output any] struct{}

func (*userAgentMiddleware[Input, Output]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, Output]) (*handlerOutput[Output], error) {
	if input.Options.UserAgent != nil && *input.Options.UserAgent != "" {
		input.ServerRequest.Header.SetUserAgent(*input.Options.UserAgent)
	}

	return next.Handle(ctx, input)
}

type resolveEndpointMiddleware[Input any, Output any] struct{}

func (*resolveEndpointMiddleware[Input, Output]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, Output]) (*handlerOutput[Output], error) {
	params := EndpointParameters{
		Host:         input.Options.EndpointHost,
		UseSSL:       input.Options.UseSSL,
		UsePathStyle: input.Options.UsePathStyle,
	}

	if v, ok := interface{}(input.CallInput).(RequiredBucketKeyInterface); ok {
		params.Bucket = v.GetBucket()
		params.Key = v.GetKey()
	} else if v, ok := interface{}(input.CallInput).(RequiredBucketInterface); ok {
		params.Bucket = v.GetBucket()
	}

	endpoint, err := input.Options.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve endpoint: %v", err)
	}

	input.ServerRequest.SetRequestURI(endpoint.URL)

	return next.Handle(ctx, input)
}

type signerMiddleware[Input any, Output any] struct{}

func (*signerMiddleware[Input, Output]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, Output]) (*handlerOutput[Output], error) {
	if _, _, err := input.Options.Signer.Sign(&input.ServerRequest, input.Options.Credentials, input.Options.SiginingRegion, time.Now()); err != nil {
		return nil, fmt.Errorf("cannot sign the request: %v", err)
	}

	return next.Handle(ctx, input)
}

type transportMiddleware[
	Input HTTPRequestMarshaler,
	OutputBase any,
	OutputPtr interface {
		HTTPResponseUnmarshaler
		*OutputBase
	},
] struct{}

func (*transportMiddleware[Input, OutputBase, OutputPtr]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, OutputPtr]) (*handlerOutput[OutputPtr], error) {
	if err := input.CallInput.MarshalHTTP(&input.ServerRequest); err != nil {
		return nil, fmt.Errorf("HTTP marshaling error: %v", err)
	}

	output, err := next.Handle(ctx, input)
	if err != nil {
		return output, err
	}

	var callOutputBase OutputBase
	callOutputPtr := OutputPtr(&callOutputBase)

	// Do not wrap error since an unexpected HTTP status code can make
	// UnmarshalHTTP to return a server-side error.
	if err := callOutputPtr.UnmarshalHTTP(output.ServerResponse); err != nil {
		return output, err
	}

	output.CallOutputV3 = callOutputPtr
	return output, nil
}

type requiredInputMiddleware[Input any, Output any] struct{}

func (*requiredInputMiddleware[Input, Output]) Middleware(ctx context.Context, input *handlerInput[Input], next Handler[Input, Output]) (*handlerOutput[Output], error) {
	if v, ok := interface{}(input.CallInput).(RequiredBucketKeyInterface); ok {
		if v.GetBucket() == "" {
			return nil, errors.New("bucket is mandatory")
		}

		if v.GetKey() == "" {
			return nil, errors.New("object key is mandatory")
		}
	} else if v, ok := interface{}(input.CallInput).(RequiredBucketInterface); ok {
		if v.GetBucket() == "" {
			return nil, errors.New("bucket is mandatory")
		}
	}

	return next.Handle(ctx, input)
}
