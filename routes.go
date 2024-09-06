package main

import (
	"constellation/handlers"
	middlewares "constellation/middleware"
	"net/http"
)

func apiRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	apiURL := "/clusters/v1"
	mux.HandleFunc("GET "+apiURL+"/ping", middlewares.BundleMiddlewares(handlers.Ping))

	return mux
}
