package handlers

import (
	"constellation/clusters/k8s"
	"log/slog"
	"net/http"
)

func GetSecrets(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	resource := k8s.NewCoreAPIResource()
	namespace := r.PathValue("namespace")
	secrets, err := resource.Secrets(namespace)
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusInternalServerError, "error", nil, nil)
		return
	}
	var secretsData []map[string]interface{}
	for _, secret := range secrets.Items {
		secretsData = append(secretsData, map[string]interface{}{
			"name": secret.Name,
			"type": secret.Type,
			"data": secret.Data,
		})
	}
	crw.response(http.StatusOK, "ok", map[string]interface{}{
		"number_of_secrets": len(secrets.Items),
		"secrets":           secretsData,
	}, nil)
}

func GetPods(w http.ResponseWriter, r *http.Request) {}

func GetConfigMaps(w http.ResponseWriter, r *http.Request) {}
