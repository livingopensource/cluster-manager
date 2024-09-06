package postgres

import "constellation/clusters"

type PostgresImpl struct {
	kubeconfig string
}

func NewCluster(config string) *PostgresImpl {
	return &PostgresImpl{
		kubeconfig: config,
	}
}

func (c *PostgresImpl) Create(resource clusters.ClusterResource) error {
	return nil
}

func (c *PostgresImpl) Delete(id string) error {
	return nil
}

func (c *PostgresImpl) Update(id string, resource clusters.ClusterResource) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (c *PostgresImpl) Find(id string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (c *PostgresImpl) FindAll() ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
