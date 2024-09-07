package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantLogs string
	}{
		{
			name:     "Root path",
			path:     "/",
			wantLogs: "/",
		},
		{
			name:     "API path",
			path:     "/api/v1/users",
			wantLogs: "/api/v1/users",
		},
		{
			name:     "Path with query parameters",
			path:     "/search?q=test",
			wantLogs: "/search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, nil))
			slog.SetDefault(logger)

			handler := logMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if !strings.Contains(buf.String(), tt.wantLogs) {
				t.Errorf("Expected logs to contain %q, but got %q", tt.wantLogs, buf.String())
			}
		})
	}
}

func TestLogMiddlewareCallsNextHandler(t *testing.T) {
	nextHandlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
	})

	handler := logMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if !nextHandlerCalled {
		t.Error("Next handler was not called")
	}
}
