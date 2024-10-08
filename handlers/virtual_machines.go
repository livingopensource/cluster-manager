package handlers

import (
	"constellation/clusters"
	"constellation/clusters/vm"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"k8s.io/apimachinery/pkg/api/errors"
)

func CreateVirtualMachineInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
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
	err = clusters.CreateResource(vmInstance, clusterResource)
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
	crw.response(http.StatusCreated, "VirtualMachine instance created", nil, nil)
}

func GetAllVirtualMachineInstances(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
	namespace := r.PathValue("namespace")
	params := r.URL.Query()
	instances, err := clusters.FindAllResources(vmInstance, clusters.ClusterResource{
		Namespace: namespace,
		HTTP: clusters.HTTP{
			QueryParams: params,
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
	crw.response(http.StatusOK, "VirtualMachine instances fetched", instances, nil)
}

func GetVirtualMachineInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	params := r.URL.Query()
	instance, err := clusters.FindResource(vmInstance, clusters.ClusterResource{
		Namespace: namespace,
		Compute: clusters.Compute{
			Name: name,
		},
		HTTP: clusters.HTTP{
			QueryParams: params,
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
	crw.response(http.StatusOK, "VirtualMachine instance fetched", instance, nil)
}

func PatchVirtualMachineInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
	namespace := r.PathValue("namespace")
	var clusterResource clusters.ClusterResource
	err := json.NewDecoder(r.Body).Decode(&clusterResource)
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}
	clusterResource.Namespace = namespace
	response, err := vmInstance.Patch(clusterResource)
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}
	crw.response(http.StatusAccepted, "ok", response, nil)
}

func DeleteVirtualMachineInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	err := clusters.DeleteResource(vmInstance, clusters.ClusterResource{
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
	crw.response(http.StatusOK, "VirtualMachine instance deleted", nil, nil)
}

func WatchVirtualMachineInstances(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
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

	watcher, err := vmInstance.Watch(clusters.ClusterResource{
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

func VNCVirtualMachineInstance(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	vmInstance := vm.NewCluster()
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")

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

	stream, err := vmInstance.VNC(clusters.ClusterResource{
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

	defer stream.AsConn().Close()

	vmiConn := stream.AsConn()
	// Start copying data between WebSocket and VMI console
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Error reading from websocket: " + err.Error())
				break
			}
			_, err = vmiConn.Write(message)
			if err != nil {
				slog.Error("Error writing to VMI console: n" + err.Error())
				break
			}
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := vmiConn.Read(buf)
		if err != nil {
			slog.Error("Error reading from VMI console: " + err.Error())
			break
		}
		err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			slog.Error("Error writing to websocket: " + err.Error())
			break
		}
	}

}
