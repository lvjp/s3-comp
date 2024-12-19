package client

import (
	"context"
	"net/http"
	"net/url"
	"slices"
)

type HTTPMarshaler interface {
	MarshalHTTP(context.Context, *http.Request) error
}

type HTTPUnmarshaler interface {
	UnmarshalHTTP(context.Context, *http.Response) error
}

type MiddlewareContext struct {
	context.Context

	Bucket *string

	HTTPRequest  *http.Request
	HTTPResponse *http.Response

	S3Input  HTTPMarshaler
	S3Output HTTPUnmarshaler
}

func NewMiddlewareContext(ctx context.Context, bucket *string, s3Input HTTPMarshaler, s3Output HTTPUnmarshaler) *MiddlewareContext {
	return &MiddlewareContext{
		Context:  ctx,
		Bucket:   bucket,
		S3Input:  s3Input,
		S3Output: s3Output,
	}
}

type Handler interface {
	Handle(*MiddlewareContext) error
}

type HandlerFunc func(*MiddlewareContext) error

func (f HandlerFunc) Handle(ctx *MiddlewareContext) error {
	return f(ctx)
}

type Middleware interface {
	HandleMiddleware(next Handler) Handler
}

type MiddlewareFunc func(next Handler) Handler

func (mw MiddlewareFunc) HandleMiddleware(next Handler) Handler {
	return mw(next)
}

func DecorateHandler(h Handler, mw ...MiddlewareFunc) Handler {
	for i := range slices.Backward(mw) {
		h = mw[i].HandleMiddleware(h)
	}

	return h
}

func initHTTPRequestMiddleware(next Handler) Handler {
	return HandlerFunc(func(ctx *MiddlewareContext) error {
		ctx.HTTPRequest = &http.Request{
			URL:    new(url.URL),
			Header: http.Header{},

			// The value -1 indicates that the length is unknown.
			ContentLength: -1,
			Body:          http.NoBody,
		}

		return next.Handle(ctx)
	})
}

func mandatoryBucketMiddleware(next Handler) Handler {
	return HandlerFunc(func(ctx *MiddlewareContext) error {
		if ctx.Bucket == nil || len(*ctx.Bucket) == 0 {
			return NewSDKErrorBucketIsMandatory(ctx)
		}

		return next.Handle(ctx)
	})
}

func httpMarshalerMiddleware(next Handler) Handler {
	return HandlerFunc(func(ctx *MiddlewareContext) error {
		if err := ctx.S3Input.MarshalHTTP(ctx, ctx.HTTPRequest); err != nil {
			return err
		}

		return next.Handle(ctx)
	})
}

func httpUnmarshalerMiddleware(next Handler) Handler {
	return HandlerFunc(func(ctx *MiddlewareContext) error {
		err := next.Handle(ctx)
		if err != nil {
			return err
		}

		if err := ctx.S3Output.UnmarshalHTTP(ctx, ctx.HTTPResponse); err != nil {
			return err
		}

		return nil
	})
}
