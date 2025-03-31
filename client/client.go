package client

import (
	"context"

	chain_of_responsibility "github.com/lvjp/s3-comp/client/pkg/design-patterns/chain-of-responsibility"

	"github.com/valyala/fasthttp"
)

type Client struct {
	options Options
}

func New(options *Options, optFns ...func(*Options)) (*Client, error) {
	c := &Client{}

	if options != nil {
		c.options = *options
	}

	for _, fn := range optFns {
		fn(&c.options)
	}

	c.options.setDefaults()
	if err := c.options.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

type Metadata struct {
	Request  *fasthttp.Request
	Response *fasthttp.Response
}

func PerformCall[
	Input HTTPRequestMarshaler,
	OutputPtr interface {
		HTTPResponseUnmarshaler
		*OutputBase
	},
	OutputBase any,
](ctx context.Context, c *Client, input Input, optFns ...func(*Options)) (OutputPtr, *Metadata, error) {
	in := &handlerInput[Input]{
		Options:   c.options.With(optFns...),
		CallInput: input,
	}

	chain := chain_of_responsibility.NewChain(
		&httpRequesterHandler[Input, OutputPtr]{},
		&errorMiddleware[Input, OutputPtr]{},
		&configValidationMiddleware[Input, OutputPtr]{},
		&requiredInputMiddleware[Input, OutputPtr]{},
		&userAgentMiddleware[Input, OutputPtr]{},
		&resolveEndpointMiddleware[Input, OutputPtr]{},
		&transportMiddleware[Input, OutputBase, OutputPtr]{},
		&signerMiddleware[Input, OutputPtr]{},
	)

	out, err := chain.Handle(ctx, in)

	metadata := &Metadata{
		Request: &in.ServerRequest,
	}

	if out != nil {
		metadata.Response = out.ServerResponse
	}

	if err != nil {
		return nil, metadata, err
	}

	return out.CallOutputV3, metadata, nil
}
