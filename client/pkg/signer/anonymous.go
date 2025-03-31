package signer

import (
	"time"

	"github.com/valyala/fasthttp"
)

func NewAnonymousSigner() Signer {
	return &anonymousSigner{}
}

type anonymousSigner struct {
}

func (*anonymousSigner) Sign(*fasthttp.Request, *Credentials, string, time.Time) (string, string, error) {
	return "", "", nil
}
