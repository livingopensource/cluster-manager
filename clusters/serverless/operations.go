package serverless

import (
	"constellation/clusters"
	"errors"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/watch"
)

type Serverless struct {
	kubeconfig string
}

func NewCluster() *Serverless {
	return &Serverless{
		kubeconfig: viper.GetString("virtual_machines.kubeconfig"),
	}
}

func (c *Serverless) Create(resource clusters.ClusterResource) error {
	return errors.New("not implemented")
}

func (c *Serverless) Delete(resource clusters.ClusterResource) error {
	return errors.New("not implemented")
}

func (c *Serverless) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *Serverless) Watch(resource clusters.ClusterResource) (watch.Interface, error) {
	return nil, errors.New("not implemented")
}

func (c *Serverless) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *Serverless) FindAll(resource clusters.ClusterResource) ([]map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *Serverless) Update(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
