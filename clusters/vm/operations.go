package vm

import (
	"constellation/clusters"
	"errors"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/watch"
)

type VirtualMachine struct {
	kubeconfig string
}

func NewCluster() *VirtualMachine {
	return &VirtualMachine{
		kubeconfig: viper.GetString("virtual_machines.kubeconfig"),
	}
}

func (c *VirtualMachine) Create(resource clusters.ClusterResource) error {
	return errors.New("not implemented")
}

func (c *VirtualMachine) Delete(resource clusters.ClusterResource) error {
	return errors.New("not implemented")
}

func (c *VirtualMachine) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *VirtualMachine) Watch(resource clusters.ClusterResource) (watch.Interface, error) {
	return nil, errors.New("not implemented")
}

func (c *VirtualMachine) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *VirtualMachine) FindAll(resource clusters.ClusterResource) ([]map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *VirtualMachine) Update(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
