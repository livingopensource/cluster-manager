package handlers

import (
	"constellation/clusters"
	"constellation/clusters/serverless"
	"constellation/clusters/vm"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"k8s.io/apimachinery/pkg/api/errors"
)

func CreateServerlessInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	serverlessInstance := serverless.NewCluster()
	namespace := r.PathValue("namespace")
	req, err := io.ReadAll(r.Body)
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	var clusterResource clusters.ClusterResource
	err = json.Unmarshal(req, &clusterResource)
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	clusterResource.Namespace = namespace
	err = clusters.CreateResource(serverlessInstance, clusterResource)
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusCreated, "Serveless instance created", nil, nil)
}

func GetAllServerlessInstances(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	serverlessInstance := serverless.NewCluster()
	namespace := r.PathValue("namespace")
	instances, err := clusters.FindAllResources(serverlessInstance, clusters.ClusterResource{
		Namespace: namespace,
	})
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "Serveless instances fetched", instances, nil)
}

func GetServerlessInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	serverlessInstance := serverless.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	instance, err := clusters.FindResource(serverlessInstance, clusters.ClusterResource{
		Namespace: namespace,
		Compute: clusters.Compute{
			Name: name,
		},
	})
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "Serveless instance fetched", instance, nil)
}

func DeleteServerlessInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	serverlessInstance := serverless.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	err := clusters.DeleteResource(serverlessInstance, clusters.ClusterResource{
		Namespace: namespace,
		Compute: clusters.Compute{
			Name: name,
		},
	})
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "Serveless instance deleted", nil, nil)
}

func WatchServerlessInstances(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	serverlessInstance := serverless.NewCluster()
	namespace := r.PathValue("namespace")
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	watcher, err := serverlessInstance.Watch(clusters.ClusterResource{
		Namespace: namespace,
	})
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	defer watcher.Stop()

	// Stream events to the websocket
	for event := range watcher.ResultChan() {
		jsonData, err := json.Marshal(event)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			slog.Error(err.Error())
			break
		}
	}
}
