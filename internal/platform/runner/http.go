package runner

import (
	"fmt"
	"log"
	"log/slog"

	httpserver "github.com/project-820/transactions/internal/platform/http"
)

var _ Runner = (*HTTPRunner)(nil)

type HTTPRunner struct {
	addr   string
	server *httpserver.HTTP
}

func NewHTTPRunner(addr string, server *httpserver.HTTP) *HTTPRunner {
	return &HTTPRunner{addr: addr, server: server}
}

func (h *HTTPRunner) Start() error {
	go func() {
		if err := h.server.Serve(); err != nil {
			log.Fatalf("failed to serve http server: %v", err)
		}
	}()
	slog.Info(fmt.Sprintf("serving http on http://localhost%s", h.addr))
	return nil
}

func (h *HTTPRunner) Stop() error {
	if err := h.server.Stop(); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}
	slog.Info("HTTP server stopped")
	return nil
}
