package k8s

import (
	"context"

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	kubeconfig string
}

func NewCoreAPIResource() *Resource {
	return &Resource{
		kubeconfig: viper.GetString("postgres_cluster.kubeconfig"),
	}
}

func (r Resource) Pods(ns string) (*v1.PodList, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
}

func (r Resource) Secrets(ns string) (*v1.SecretList, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().Secrets(ns).List(context.Background(), metav1.ListOptions{})
}

func (r Resource) ConfigMaps(ns string) (*v1.ConfigMapList, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().ConfigMaps(ns).List(context.Background(), metav1.ListOptions{})
}
