package handlers

import (
	"constellation/clusters"
	"constellation/clusters/postgres"
	"log/slog"
	"net/http"
)

func CreatePostgresInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	pgInstance := postgres.NewCluster()

	err := clusters.CreateResource(pgInstance, clusters.ClusterResource{})
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
	}
	crw.response(http.StatusCreated, "Postgres instance created", nil, nil)
}

func GetAllPostgresInstances(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	pgInstance := postgres.NewCluster()
	namespace := r.PathValue("namespace")
	instances, err := clusters.FindAllResources(pgInstance, clusters.ClusterResource{
		Namespace: namespace,
	})
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
	}
	crw.response(http.StatusOK, "Postgres instances fetched", instances, nil)
}

func GetPostgresInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	pgInstance := postgres.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	instance, err := clusters.FindResource(pgInstance, clusters.ClusterResource{
		Namespace: namespace,
		Compute: clusters.Compute{
			Name: name,
		},
	})
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
	}
	crw.response(http.StatusOK, "Postgres instance fetched", instance, nil)
}

func DeletePostgresInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	pgInstance := postgres.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	err := clusters.DeleteResource(pgInstance, clusters.ClusterResource{
		Namespace: namespace,
		Compute: clusters.Compute{
			Name: name,
		},
	})
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
	}
	crw.response(http.StatusOK, "Postgres instance deleted", nil, nil)
}
