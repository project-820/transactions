package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type HTTPParams struct {
	Address string
	Handler http.Handler

	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

type HTTP struct {
	server          *http.Server
	listener        net.Listener
	shutdownTimeout time.Duration
}

func NewHTTPServer(params HTTPParams) (*HTTP, error) {
	if params.Address == "" {
		return nil, fmt.Errorf("http address is empty")
	}

	if params.Handler == nil {
		return nil, fmt.Errorf("http handler is nil")
	}

	ln, err := net.Listen("tcp", params.Address)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %w", params.Address, err)
	}

	st := params.ShutdownTimeout
	if st == 0 {
		st = 10 * time.Second
	}

	mhb := params.MaxHeaderBytes
	if mhb <= 0 {
		mhb = 1 << 20 // 1MiB
	}

	srv := &http.Server{
		Addr:           params.Address,
		Handler:        params.Handler,
		ReadTimeout:    params.ReadTimeout,
		WriteTimeout:   params.WriteTimeout,
		IdleTimeout:    params.IdleTimeout,
		MaxHeaderBytes: mhb,
	}

	return &HTTP{server: srv, listener: ln, shutdownTimeout: st}, nil
}

func (h *HTTP) Serve() error {
	slog.Info(fmt.Sprintf("http server listening on %s", h.listener.Addr().String()))

	err := h.server.Serve(h.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (h *HTTP) Stop() error {
	slog.Info("shutting down http server...")

	ctx, cancel := context.WithTimeout(context.Background(), h.shutdownTimeout)
	defer cancel()

	return h.server.Shutdown(ctx)
}
