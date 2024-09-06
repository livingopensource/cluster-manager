package middlewares

import (
	"log/slog"
	"net/http"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
