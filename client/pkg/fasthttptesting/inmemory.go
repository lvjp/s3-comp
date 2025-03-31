package fasthttptesting

import (
	"net"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func NewInmemoryTester(handler func(ctx *fasthttp.RequestCtx)) InmemoryTester {
	ret := &inmemoryTester{
		ln: fasthttputil.NewInmemoryListener(),
		server: &fasthttp.Server{
			Handler: handler,
		},
	}

	ret.serverWg.Add(1)
	go func() {
		defer ret.serverWg.Done()
		ret.serverErr = ret.server.Serve(ret.ln)
	}()

	return ret
}

type InmemoryTester interface {
	Close()
	Client() *fasthttp.Client
}

type inmemoryTester struct {
	ln        *fasthttputil.InmemoryListener
	server    *fasthttp.Server
	serverErr error
	serverWg  sync.WaitGroup
}

func (in *inmemoryTester) Close() {
	if in.server != nil {
		err := in.server.Shutdown()
		if err != nil {
			panic("inmemoryTester: cannot shutdown server: " + err.Error())
		}
		in.server = nil
		in.ln = nil

		in.serverWg.Wait()
		if in.serverErr != nil {
			panic("inmemoryTester: server error: " + in.serverErr.Error())
		}
	}
}

func (in *inmemoryTester) Client() *fasthttp.Client {
	return &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		Dial: func(addr string) (net.Conn, error) {
			return in.ln.Dial()
		},
	}
}
