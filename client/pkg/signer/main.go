package signer

import (
	"time"

	"github.com/valyala/fasthttp"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type Signer interface {
	Sign(req *fasthttp.Request, credentials *Credentials, region string, now time.Time) (canonicalRequest, stringToSign string, err error)
}
