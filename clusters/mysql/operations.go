package vm

import (
	"constellation/clusters"
	"errors"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/watch"
)

type MySQL struct {
	kubeconfig string
}

func NewCluster() *MySQL {
	return &MySQL{
		kubeconfig: viper.GetString("virtual_machines.kubeconfig"),
	}
}

func (c *MySQL) Create(resource clusters.ClusterResource) error {
	return errors.New("not implemented")
}

func (c *MySQL) Delete(resource clusters.ClusterResource) error {
	return errors.New("not implemented")
}

func (c *MySQL) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *MySQL) Watch(resource clusters.ClusterResource) (watch.Interface, error) {
	return nil, errors.New("not implemented")
}

func (c *MySQL) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *MySQL) FindAll(resource clusters.ClusterResource) ([]map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *MySQL) Update(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
