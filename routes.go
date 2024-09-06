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

	db := apiURL + "/databases"
	postgres := db + "/postgres"
	ns := postgres + "/{namespace}"
	mux.HandleFunc("POST "+ns, middlewares.BundleMiddlewares(handlers.CreatePostgresInstance))
	mux.HandleFunc("GET "+ns, middlewares.BundleMiddlewares(handlers.GetAllPostgresInstances))
	mux.HandleFunc("GET "+ns+"/{name}", middlewares.BundleMiddlewares(handlers.GetPostgresInstance))
	mux.HandleFunc("DELETE "+ns+"/{name}", middlewares.BundleMiddlewares(handlers.DeletePostgresInstance))

	return mux
}
