package clusters

import "k8s.io/apimachinery/pkg/watch"

type ClusterImpl interface {
	Create(resource ClusterResource) error
	Update(resource ClusterResource) (map[string]interface{}, error)
	Delete(resource ClusterResource) error
	Find(resource ClusterResource) (map[string]interface{}, error)
	FindAll(resource ClusterResource) ([]map[string]interface{}, error)
	Watch(resource ClusterResource) (watch.Interface, error)
	Patch(resource ClusterResource) (map[string]interface{}, error)
}

func CreateResource(r ClusterImpl, resource ClusterResource) error {
	return r.Create(resource)
}

func UpdateResource(r ClusterImpl, resource ClusterResource) (map[string]interface{}, error) {
	return r.Update(resource)
}

func DeleteResource(r ClusterImpl, resource ClusterResource) error {
	return r.Delete(resource)
}

func FindResource(r ClusterImpl, resource ClusterResource) (map[string]interface{}, error) {
	return r.Find(resource)
}

func FindAllResources(r ClusterImpl, resource ClusterResource) ([]map[string]interface{}, error) {
	return r.FindAll(resource)
}

func Watch(r ClusterImpl, resource ClusterResource) (watch.Interface, error) {
	return r.Watch(resource)
}

func Patch(r ClusterImpl, resource ClusterResource) (map[string]interface{}, error) {
	return r.Patch(resource)
}
