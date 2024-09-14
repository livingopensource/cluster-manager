package postgres

import (
	"constellation/clusters"
	"errors"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

type PostgresImpl struct {
	kubeconfig string
}

func NewCluster() *PostgresImpl {
	return &PostgresImpl{
		kubeconfig: viper.GetString("postgres_cluster.kubeconfig"),
	}
}

func (c *PostgresImpl) Create(resource clusters.ClusterResource) error {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "postgresql.cnpg.io/v1",
			"kind":       "Cluster",
			"metadata": map[string]interface{}{
				"name": resource.Compute.Name,
			},
			"spec": map[string]interface{}{
				"instances": resource.Compute.Instances,
				"bootstrap": map[string]interface{}{
					"initdb": map[string]interface{}{
						"database": resource.Compute.Name,
						"owner":    resource.Account.Name,
					},
				},
				"resources": map[string]interface{}{
					"limits": map[string]interface{}{
						"cpu":    resource.Compute.CPU,
						"memory": resource.Compute.RAM,
					},
				},
				"storage": map[string]interface{}{
					"size": resource.Compute.Storage,
				},
			},
		},
	}
	_, err := clusters.CreateResourceSchema(obj, c.kubeconfig, resource.Namespace)
	return err
}

func (c *PostgresImpl) Delete(resource clusters.ClusterResource) error {
	return clusters.DeleteResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace)
}

func (c *PostgresImpl) Update(resource clusters.ClusterResource) (map[string]interface{}, error) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "postgresql.cnpg.io/v1",
			"kind":       "Cluster",
			"metadata": map[string]interface{}{
				"name": resource.Compute.Name,
			},
			"spec": map[string]interface{}{
				"instances": resource.Compute.Instances,
				"storage": map[string]interface{}{
					"size": resource.Compute.Storage,
				},
			},
		},
	}
	response, err := clusters.UpdateResourceSchema(obj, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	return response.Object, err
}

func (c *PostgresImpl) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *PostgresImpl) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	response, err := clusters.GetResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	return response.Object, err
}

func (c *PostgresImpl) FindAll(resource clusters.ClusterResource) ([]map[string]interface{}, error) {
	response, err := clusters.ListResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(response.Items))
	for i, item := range response.Items {
		result[i] = item.Object
	}
	return result, nil
}

func (c *PostgresImpl) Watch(resource clusters.ClusterResource) (watch.Interface, error) {
	response, err := clusters.WatchResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	return response, err
}
