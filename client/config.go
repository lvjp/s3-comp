package client

import (
	"net/http"
)

type Config struct {
	HTTPClient HTTPClient
}

func (cfg *Config) SetDefaults() {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
