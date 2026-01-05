package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewDefaultClient(timeout time.Duration) Doer {
	return &http.Client{Timeout: timeout}
}

func DoJSON(ctx context.Context, d Doer, method, url string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := d.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	return b, nil
}
