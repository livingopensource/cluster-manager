package clusters

import (
	"constellation/clusters/k8s"
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateResourceSchema(resource *unstructured.Unstructured, config, namespace, service string) (*unstructured.Unstructured, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, err
	}
	gvk := resource.GetObjectKind().GroupVersionKind()
	rm := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	dyn, err := k8s.DynamicClientSet(config, service)
	if err != nil {
		return nil, err
	}
	if mapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = dyn.Resource(mapping.Resource)
	} else {
		ri = dyn.Resource(mapping.Resource).Namespace(namespace)
	}

	return ri.Create(context.TODO(), resource, metav1.CreateOptions{})
}

func UpdateResourceSchema(resource *unstructured.Unstructured, config, namespace, service string) (*unstructured.Unstructured, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, err
	}
	gvk := resource.GetObjectKind().GroupVersionKind()
	rm := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	dyn, err := k8s.DynamicClientSet(config, service)
	if err != nil {
		return nil, err
	}
	if mapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = dyn.Resource(mapping.Resource)
	} else {
		ri = dyn.Resource(mapping.Resource).Namespace(namespace)
	}

	return ri.Update(context.TODO(), resource, metav1.UpdateOptions{})
}

func GetResourceSchema(gvk schema.GroupVersionKind, name, config, namespace, service string) (*unstructured.Unstructured, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, err
	}

	rm := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	dyn, err := k8s.DynamicClientSet(config, service)
	if err != nil {
		return nil, err
	}
	if mapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = dyn.Resource(mapping.Resource)
	} else {
		ri = dyn.Resource(mapping.Resource).Namespace(namespace)
	}

	return ri.Get(context.TODO(), name, metav1.GetOptions{})
}

func ListResourceSchema(gvk schema.GroupVersionKind, config, namespace, service string) (*unstructured.UnstructuredList, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, err
	}

	rm := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	dyn, err := k8s.DynamicClientSet(config, service)
	if err != nil {
		return nil, err
	}
	if mapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = dyn.Resource(mapping.Resource)
	} else {
		ri = dyn.Resource(mapping.Resource).Namespace(namespace)
	}

	return ri.List(context.TODO(), metav1.ListOptions{})
}

func DeleteResourceSchema(gvk schema.GroupVersionKind, name, config, namespace, service string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return err
	}

	rm := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	var ri dynamic.ResourceInterface
	dyn, err := k8s.DynamicClientSet(config, service)
	if err != nil {
		return err
	}
	if mapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = dyn.Resource(mapping.Resource)
	} else {
		ri = dyn.Resource(mapping.Resource).Namespace(namespace)
	}

	return ri.Delete(context.TODO(), name, metav1.DeleteOptions{})
}
