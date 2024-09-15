package vm

import (
	"constellation/clusters"
	"errors"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	return clusters.DeleteResourceSchema(schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace)
}

func (c *VirtualMachine) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c *VirtualMachine) Watch(resource clusters.ClusterResource) (watch.Interface, error) {
	response, err := clusters.WatchResourceSchema(schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *VirtualMachine) Find(resource clusters.ClusterResource) (map[string]interface{}, error) {
	gvk := schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}
	if resource.HTTP.QueryParams.Get("state") == "up" {
		gvk = schema.GroupVersionKind{
			Group:   "kubevirt.io",
			Version: "v1",
			Kind:    "VirtualMachineInstance",
		}
	}
	response, err := clusters.GetResourceSchema(gvk, resource.Compute.Name, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	return response.Object, nil
}

func (c *VirtualMachine) FindAll(resource clusters.ClusterResource) ([]map[string]interface{}, error) {
	gvk := schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}
	if resource.HTTP.QueryParams.Get("state") == "up" {
		gvk = schema.GroupVersionKind{
			Group:   "kubevirt.io",
			Version: "v1",
			Kind:    "VirtualMachineInstance",
		}
	}
	response, err := clusters.ListResourceSchema(gvk, c.kubeconfig, resource.Namespace)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(response.Items))
	for i, item := range response.Items {
		result[i] = item.Object
	}
	return result, nil
}

func (c *VirtualMachine) Update(resource clusters.ClusterResource) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
