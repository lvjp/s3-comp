package client

import (
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

func setQuery(args *fasthttp.Args, name string, value *string) {
	if value == nil {
		return
	}

	args.Set(name, *value)
}

func setHeader(requestHeader *fasthttp.RequestHeader, key string, value *string) {
	if value == nil {
		return
	}

	requestHeader.Set(key, *value)
}

func setHeaderOrTrailer(requestHeader *fasthttp.RequestHeader, key string, header, trailer *string) {
	haveValue := header != nil
	haveTrailer := trailer != nil

	switch {
	case haveValue && haveTrailer:
		panic("cannot set both header and trailer for " + strconv.Quote(key))
	case haveValue && !haveTrailer:
		requestHeader.Set(key, *header)
	case !haveValue && haveTrailer:
		if actual := requestHeader.Peek(HeaderXAmzTrailer); len(actual) > 0 {
			panic(fmt.Sprintf("trailer already set: %q", actual))
		}

		requestHeader.Set(HeaderXAmzTrailer, key)
		requestHeader.Set(key, *trailer)
	default:
		// !haveValue && !haveTrailer is a no-op
	}
}

func extractHeader(responseHeader *fasthttp.ResponseHeader, key string) *string {
	value := responseHeader.Peek(key)
	if value == nil {
		return nil
	}

	str := string(value)
	return &str
}
