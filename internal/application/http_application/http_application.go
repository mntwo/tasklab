package http_application

import (
	"context"
	"net/http"

	"github.com/mntwo/tasklab/internal/application"
	"github.com/mntwo/tasklab/internal/config"
)

var _ application.Application = (*HttpApplication)(nil)

type HttpApplication struct {
	Name    string
	Server  *http.Server
	Handler http.Handler
}

func New(name string, opts ...Option) *HttpApplication {
	cfg := defaultConfig(name)

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &HttpApplication{
		Name:    name,
		Server:  cfg.server,
		Handler: cfg.handler,
	}
}

func defaultConfig(name string) *optconfig {
	c := config.GetHttpServer(name)
	srv := &http.Server{
		Addr:         c.Addr,
		WriteTimeout: c.WriteTimeout,
		ReadTimeout:  c.ReadTimeout,
		IdleTimeout:  c.IdleTimeout,
	}
	handler := http.DefaultServeMux
	return &optconfig{
		server:  srv,
		handler: handler,
	}
}

func (a *HttpApplication) Start() error {
	a.Server.Handler = a.Handler
	if err := a.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a *HttpApplication) Stop() error {
	c := config.GetHttpServer(a.Name)
	ctx, cancel := context.WithTimeout(context.Background(), c.CloseTimeout)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (a *HttpApplication) GetName() string {
	return a.Name
}
