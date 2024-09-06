package clusters

type ClusterImpl interface {
	Create(resource ClusterResource) error
	Update(id string, resource ClusterResource) (map[string]interface{}, error)
	Delete(id string) error
	Find(id string) (map[string]interface{}, error)
	FindAll() ([]map[string]interface{}, error)
}

func CreateResource(r ClusterImpl, resource ClusterResource) error {
	return r.Create(resource)
}

func UpdateResource(r ClusterImpl, resource ClusterResource) (map[string]interface{}, error) {
	return r.Update(resource.ID, resource)
}

func DeleteResource(r ClusterImpl, resource ClusterResource) error {
	return r.Delete(resource.ID)
}

func FindResource(r ClusterImpl, resource ClusterResource) (map[string]interface{}, error) {
	return r.Find(resource.ID)
}

func FindAllResources(r ClusterImpl, resource ClusterResource) ([]map[string]interface{}, error) {
	return r.FindAll()
}
