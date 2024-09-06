package postgres

import (
	"constellation/clusters"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
				"storage": map[string]interface{}{
					"size": resource.Compute.Storage,
				},
			},
		},
	}
	_, err := clusters.CreateResourceSchema(obj, c.kubeconfig, resource.Namespace, "postgres")
	return err
}

func (c *PostgresImpl) Delete(resource clusters.ClusterResource) error {
	return clusters.DeleteResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace, "postgres")
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
	response, err := clusters.UpdateResourceSchema(obj, c.kubeconfig, resource.Namespace, "postgres")
	return response.Object, err
}

func (c *PostgresImpl) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	response, err := clusters.GetResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace, "postgres")
	return response.Object, err
}

func (c *PostgresImpl) FindAll(resource clusters.ClusterResource) (map[string]interface{}, error) {
	response, err := clusters.ListResourceSchema(schema.GroupVersionKind{
		Group:   "postgresql.cnpg.io",
		Version: "v1",
		Kind:    "Cluster",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace, "postgres")
	return response.Object, err
}
