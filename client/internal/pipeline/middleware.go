package pipeline

import (
	"context"
	"net/http"
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

func NewPipeline(h Handler, mw ...MiddlewareFunc) Handler {
	for i := range slices.Backward(mw) {
		h = mw[i].HandleMiddleware(h)
	}

	return h
}
