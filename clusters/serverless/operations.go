package serverless

import (
	"constellation/clusters"
	"errors"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

type Serverless struct {
	kubeconfig string
}

func NewCluster() *Serverless {
	return &Serverless{
		kubeconfig: viper.GetString("serverless_cluster.kubeconfig"),
	}
}

func (c *Serverless) Create(resource clusters.ClusterResource) error {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "serving.knative.dev/v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": resource.Compute.Name,
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"autoscaling.knative.dev/min-scale": "0",
							"autoscaling.knative.dev/max-scale": resource.Compute.Instances,
						},
					},
					"spec": map[string]interface{}{
						"containers": map[string]interface{}{
							"image": resource.Compute.URL,
						},
					},
				},
			},
		},
	}
	_, err := clusters.CreateResourceSchema(obj, c.kubeconfig, resource.Namespace)
	return err
}

func (c *Serverless) Delete(resource clusters.ClusterResource) error {
	return clusters.DeleteResourceSchema(schema.GroupVersionKind{
		Group:   "serving.knative.dev",
		Version: "v1",
		Kind:    "Service",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace)
}

func (c *Serverless) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *Serverless) Watch(resource clusters.ClusterResource) (watch.Interface, error) {
	return clusters.WatchResourceSchema(schema.GroupVersionKind{
		Group:   "serving.knative.dev",
		Version: "v1",
		Kind:    "Service",
	}, c.kubeconfig, resource.Namespace)
}

func (c *Serverless) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	svc, err := clusters.GetResourceSchema(schema.GroupVersionKind{
		Group:   "serving.knative.dev",
		Version: "v1",
		Kind:    "Service",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	return svc.Object, nil
}

func (c *Serverless) FindAll(resource clusters.ClusterResource) ([]map[string]interface{}, error) {
	svcs, err := clusters.ListResourceSchema(schema.GroupVersionKind{
		Group:   "serving.knative.dev",
		Version: "v1",
		Kind:    "Service",
	}, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(svcs.Items))
	for i, item := range svcs.Items {
		result[i] = item.Object
	}
	return result, nil
}

func (c *Serverless) Update(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
