package client

import (
	"net/http"
	"net/url"

	"github.com/lvjp/s3-comp/client/internal/pipeline"
)

func initHTTPRequestMiddleware(next pipeline.Handler) pipeline.Handler {
	return pipeline.HandlerFunc(func(ctx *pipeline.MiddlewareContext) error {
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

func mandatoryBucketMiddleware(next pipeline.Handler) pipeline.Handler {
	return pipeline.HandlerFunc(func(ctx *pipeline.MiddlewareContext) error {
		if ctx.Bucket == nil || len(*ctx.Bucket) == 0 {
			return NewSDKErrorBucketIsMandatory(ctx)
		}

		return next.Handle(ctx)
	})
}

func httpMarshalerMiddleware(next pipeline.Handler) pipeline.Handler {
	return pipeline.HandlerFunc(func(ctx *pipeline.MiddlewareContext) error {
		if err := ctx.S3Input.MarshalHTTP(ctx, ctx.HTTPRequest); err != nil {
			return err
		}

		return next.Handle(ctx)
	})
}

func httpUnmarshalerMiddleware(next pipeline.Handler) pipeline.Handler {
	return pipeline.HandlerFunc(func(ctx *pipeline.MiddlewareContext) error {
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

func userAgentMiddleware(userAgent *string) pipeline.MiddlewareFunc {
	return func(next pipeline.Handler) pipeline.Handler {
		return pipeline.HandlerFunc(func(ctx *pipeline.MiddlewareContext) error {
			var value string
			if userAgent != nil {
				value = *userAgent
			}

			ctx.HTTPRequest.Header.Set("User-Agent", value)

			return next.Handle(ctx)
		})
	}
}

func disableDefaultTransportGzipCompression(next pipeline.Handler) pipeline.Handler {
	return pipeline.HandlerFunc(func(ctx *pipeline.MiddlewareContext) error {
		current := ctx.HTTPRequest.Header.Get("Accept-Encoding")
		if current == "" {
			ctx.HTTPRequest.Header.Set("Accept-Encoding", "")
		}

		return next.Handle(ctx)
	})
}
