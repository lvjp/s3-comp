package chain_of_responsibility

import (
	"context"
)

type Handler[Input, Output any] interface {
	Handle(ctx context.Context, input Input) (output Output, err error)
}

type HandlerFunc[Input, Output any] func(ctx context.Context, input Input) (output Output, err error)

func (fn HandlerFunc[Input, Output]) Handle(ctx context.Context, input Input) (Output, error) {
	return fn(ctx, input)
}

type Middleware[Input, Output any] interface {
	Middleware(ctx context.Context, input Input, next Handler[Input, Output]) (output Output, err error)
}

type MiddlewareFunc[Input, Output any] func(ctx context.Context, input Input, next Handler[Input, Output]) (output Output, err error)

func (fn MiddlewareFunc[Input, Output]) Middleware(ctx context.Context, input Input, next Handler[Input, Output]) (Output, error) {
	return fn(ctx, input, next)
}

type chainLink[Input, Output any] struct {
	next Handler[Input, Output]
	with Middleware[Input, Output]
}

func (chainLink chainLink[Input, Output]) Handle(ctx context.Context, input Input) (output Output, err error) {
	return chainLink.with.Middleware(ctx, input, chainLink.next)
}

func NewChain[Input, Output any](h Handler[Input, Output], with ...Middleware[Input, Output]) Handler[Input, Output] {
	for i := len(with) - 1; i >= 0; i-- {
		h = chainLink[Input, Output]{
			next: h,
			with: with[i],
		}
	}

	return h
}
