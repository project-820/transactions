package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	httpserver "github.com/project-820/transactions/internal/platform/http"
)

type HTTPConfig struct {
	Addr string
}

func NewHTTPServer(cfg *HTTPConfig, handler http.Handler) (*httpserver.HTTP, error) {
	if cfg == nil {
		return nil, fmt.Errorf("http cfg is nil")
	}
	if cfg.Addr == "" {
		return nil, fmt.Errorf("http addr is empty")
	}

	srv, err := httpserver.NewHTTPServer(httpserver.HTTPParams{
		Address:         cfg.Addr,
		Handler:         handler,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		MaxHeaderBytes:  1 << 20,
	})
	if err != nil {
		return nil, fmt.Errorf("create http server: %w", err)
	}
	return srv, nil
}
