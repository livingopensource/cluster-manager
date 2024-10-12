package vm

import (
	"constellation/clusters"
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	kv1 "kubevirt.io/api/core/v1"
	kvV1 "kubevirt.io/client-go/generated/kubevirt/clientset/versioned/typed/core/v1"
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
	name := resource.Account.Name
	passwd := resource.Account.Password
	cloudInitConfig := fmt.Sprintf(`#cloud-config
users:
  - name: %s
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: users
    home: /home/%s
    shell: /bin/bash
    lock_passwd: false
chpasswd:
  list: |
    %s:%s
  expire: False`, name, name, name, passwd)
	cloudInitBase64 := base64.StdEncoding.EncodeToString([]byte(cloudInitConfig))
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubevirt.io/v1",
			"kind":       "VirtualMachine",
			"metadata": map[string]interface{}{
				"name": resource.Compute.Name,
			},
			"spec": map[string]interface{}{
				"runStrategy": "RerunOnFailure",
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"kubevirt.io/vm": resource.Compute.Name,
						},
					},
					"spec": map[string]interface{}{
						"domain": map[string]interface{}{
							"cpu": map[string]interface{}{
								"cores": resource.Compute.CPU,
							},
							"devices": map[string]interface{}{
								"disks": []map[string]interface{}{
									{
										"name": "os-disk-" + resource.Compute.Name,
										"disk": map[string]interface{}{
											"bus": "virtio",
										},
									},
									{
										"name": "cloudinitdisk",
										"cdrom": map[string]interface{}{
											"bus": "sata",
										},
									},
								},
								"interfaces": []map[string]interface{}{
									{
										"name":       "default",
										"masquerade": map[string]interface{}{},
									},
								},
							},
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"memory": resource.Compute.RAM,
								},
							},
						},
						"networks": []map[string]interface{}{
							{
								"name": "default",
								"pod":  map[string]interface{}{},
							},
						},
						"volumes": []map[string]interface{}{
							{
								"name": "os-disk-" + resource.Compute.Name,
								"dataVolume": map[string]interface{}{
									"name": "os-volume-disk-" + resource.Compute.Name,
								},
							},
							{
								"name": "cloudinitdisk",
								"cloudInitNoCloud": map[string]interface{}{
									"userDataBase64": cloudInitBase64,
								},
							},
						},
					},
				},
				"dataVolumeTemplates": []map[string]interface{}{
					{
						"apiVersion": "cdi.kubevirt.io/v1beta1",
						"kind":       "DataVolume",
						"metadata": map[string]interface{}{
							"name": "os-volume-disk-" + resource.Compute.Name,
						},
						"spec": map[string]interface{}{
							"storage": map[string]interface{}{
								"accessModes": []string{
									"ReadWriteMany",
								},
								"resources": map[string]interface{}{
									"requests": map[string]interface{}{
										"storage": resource.Compute.Storage,
									},
								},
							},
							"source": map[string]interface{}{
								"http": map[string]interface{}{
									"url": resource.Compute.URL,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := clusters.CreateResourceSchema(obj, c.kubeconfig, resource.Namespace)
	if err != nil {
		return err
	}
	return nil
}

func (c *VirtualMachine) Delete(resource clusters.ClusterResource) error {
	return clusters.DeleteResourceSchema(schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}, resource.Compute.Name, c.kubeconfig, resource.Namespace)
}

func (c *VirtualMachine) Patch(resource clusters.ClusterResource) (map[string]interface{}, error) {
	var running bool
	var err error
	switch resource.Compute.State {
	case "on":
		running = true
	case "off":
		running = false
	default:
		err = errors.New("user did not specify wether to power on or off")
	}
	if err != nil {
		return nil, err
	}
	kubevirt, err := clusters.KubevirtResourceSchema(c.kubeconfig)
	if running {
		if err != nil {
			return map[string]interface{}{}, err
		}
		err = kubevirt.VirtualMachine(resource.Namespace).Start(context.Background(), resource.Compute.Name, &kv1.StartOptions{})
	} else {
		err = kubevirt.VirtualMachine(resource.Namespace).Stop(context.Background(), resource.Compute.Name, &kv1.StopOptions{})
	}
	return map[string]interface{}{}, err
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

func (c *VirtualMachine) VNC(resource clusters.ClusterResource) (kvV1.StreamInterface, error) {
	kubevirt, err := clusters.KubevirtResourceSchema(c.kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubevirt.VirtualMachineInstance(resource.Namespace).VNC(resource.Compute.Name)
}
