package middlewares

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func BundleMiddlewares(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	// Always run the log middleware
	m = append(m, func(next http.HandlerFunc) http.HandlerFunc {
		return logMiddleware(next).ServeHTTP
	})

	wrapped := h
	// Loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}
