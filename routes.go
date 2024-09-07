package main

import (
	"constellation/handlers"
	middlewares "constellation/middleware"
	"net/http"
)

func apiRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	ping := http.NewServeMux()

	ping.HandleFunc("GET /ping", handlers.Ping)

	// Postgres Paths
	postgres := http.NewServeMux()
	postgres.HandleFunc("POST /{namespace}/instances", handlers.CreatePostgresInstance)
	postgres.HandleFunc("GET /{namespace}/instances", handlers.GetAllPostgresInstances)
	postgres.HandleFunc("GET /{namespace}/instances/{name}", handlers.GetPostgresInstance)
	postgres.HandleFunc("DELETE /{namespace}/instances/{name}", handlers.DeletePostgresInstance)

	// MySQL paths
	mysql := http.NewServeMux()

	// Serveless app paths
	serverless := http.NewServeMux()

	// Virtual Machines
	vm := http.NewServeMux()

	// Route path subrouting
	mux.Handle("/clusters/v1/databases/postgres/", middlewares.BundleMiddlewares(http.HandlerFunc(http.StripPrefix("/clusters/v1/databases/postgres", postgres).ServeHTTP)))
	mux.Handle("/clusters/v1/", middlewares.BundleMiddlewares(http.HandlerFunc(http.StripPrefix("/clusters/v1", ping).ServeHTTP)))
	mux.Handle("/clusters/v1/databases/mysql/", middlewares.BundleMiddlewares(http.HandlerFunc(http.StripPrefix("/clusters/v1/databases/mysql", mysql).ServeHTTP)))
	mux.Handle("/clusters/v1/serverless/", middlewares.BundleMiddlewares(http.HandlerFunc(http.StripPrefix("/clusters/v1/serverless", serverless).ServeHTTP)))
	mux.Handle("/clusters/v1/virtual-machines/", middlewares.BundleMiddlewares(http.HandlerFunc(http.StripPrefix("/clusters/v1/virtual-machines", vm).ServeHTTP)))

	return mux
}
