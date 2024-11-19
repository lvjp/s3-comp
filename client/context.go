package client

import (
	"context"
)

type contextKey string

var contextKeyOperationName contextKey = "operationName"

func withOperationName(parent context.Context, operationName string) context.Context {
	return context.WithValue(parent, contextKeyOperationName, operationName)
}

func getOperationName(ctx context.Context) string {
	operationName, ok := ctx.Value(contextKeyOperationName).(string)
	if !ok {
		panic("operation name not set in context")
	}

	return operationName
}
