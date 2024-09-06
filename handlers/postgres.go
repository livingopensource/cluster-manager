package handlers

import (
	"constellation/clusters"
	"constellation/clusters/postgres"
	"net/http"
)

func CreatePostgresInstance(w http.ResponseWriter, r *http.Request) {
	pgInstance := postgres.NewCluster()
	clusters.CreateResource(pgInstance, clusters.ClusterResource{})
}
