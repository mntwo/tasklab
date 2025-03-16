package http_application

import "net/http"

type Option interface {
	apply(cfg *optconfig)
}

type option func(cfg *optconfig)

func (fn option) apply(cfg *optconfig) {
	fn(cfg)
}

type optconfig struct {
	server  *http.Server
	handler http.Handler
}

func WithHandler(handler http.Handler) Option {
	return option(func(cfg *optconfig) {
		cfg.handler = handler
	})
}
