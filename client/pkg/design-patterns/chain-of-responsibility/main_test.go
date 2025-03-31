package chain_of_responsibility

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type contextKey string

var callerKey contextKey = "the-caller"

func newMockMiddleware(v string) Middleware[string, string] {
	return MiddlewareFunc[string, string](func(ctx context.Context, input string, next Handler[string, string]) (string, error) {
		fmt.Println(">>>", v)
		fmt.Println("  ctx.caller:", ctx.Value(callerKey))
		fmt.Println("  input:", input)
		ret, err := next.Handle(context.WithValue(ctx, callerKey, v), "input-"+v)
		fmt.Println("  ret:", ret)
		fmt.Println("  err:", err)
		fmt.Println("<<<", v)

		return "ret-" + v, errors.New("error-" + v)
	})
}

func TestNewChain(t *testing.T) {
	handlerInput := "handler-input"
	handlerOutput := "handler-output"
	handlerError := errors.New("handler-error")

	handler := HandlerFunc[string, string](func(_ context.Context, input string) (string, error) {
		require.Equal(t, handlerInput, input)
		return handlerOutput, handlerError
	})

	chain := NewChain(handler)

	actualOutput, actualError := chain.Handle(context.Background(), handlerInput)
	require.Same(t, handlerError, actualError)
	require.Equal(t, handlerOutput, actualOutput)
}

func ExampleNewChain() {
	m1 := newMockMiddleware("m1")
	m2 := newMockMiddleware("m2")
	m3 := newMockMiddleware("m3")

	h := HandlerFunc[string, string](func(ctx context.Context, input string) (string, error) {
		fmt.Println(">>> handler")
		fmt.Println("  ctx.caller:", ctx.Value(callerKey))
		fmt.Println("  input:", input)
		fmt.Println("<<< handler")
		return "ret-handler", errors.New("error-handler")
	})

	chain := NewChain(h, m1, m2, m3)
	output, err := chain.Handle(context.WithValue(context.Background(), callerKey, "chain"), "input")
	fmt.Println("ret:", output)
	fmt.Println("err:", err)

	// Output:
	// >>> m1
	//   ctx.caller: chain
	//   input: input
	// >>> m2
	//   ctx.caller: m1
	//   input: input-m1
	// >>> m3
	//   ctx.caller: m2
	//   input: input-m2
	// >>> handler
	//   ctx.caller: m3
	//   input: input-m3
	// <<< handler
	//   ret: ret-handler
	//   err: error-handler
	// <<< m3
	//   ret: ret-m3
	//   err: error-m3
	// <<< m2
	//   ret: ret-m2
	//   err: error-m2
	// <<< m1
	// ret: ret-m1
	// err: error-m1
}
