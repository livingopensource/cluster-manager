package handlers

import "net/http"

func Ping(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	w.Header().Set("Content-Type", "application/json")
	crw.response(http.StatusOK, "pong", nil, nil)
}
