package middlewares

import (
	"log/slog"
	"net/http"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.URL.Path, "method", r.Method)
		next.ServeHTTP(w, r)
	})
}

func apiHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("User-Agent", "Constellation")
		w.Header().Set("X-Powered-By", "SwiftCloud")
		next.ServeHTTP(w, r)
	})
}
