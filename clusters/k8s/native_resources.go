package k8s

import (
	"context"

	"github.com/spf13/viper"
	authenticationv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	kubeconfig string
}

// NewSwiftCoreAPIResource returns a new instance of the CoreAPIResource struct.
// The variadic clusters variable is used to specify the cluster to use. Eventhough the 
// function can take multiple clusters, only the first one is used.
func NewSwiftCoreAPIResource(clusters ...string) *Resource {
	var cluster string
	if len(clusters) > 0 {
		cluster = clusters[0]
	}
	switch cluster {
	case "postgres":
		return &Resource{
			kubeconfig: viper.GetString("postgres_cluster.kubeconfig"),
		}
	case "mysql":
		return &Resource{
			kubeconfig: viper.GetString("mysql_cluster.kubeconfig"),
		}
	case "serverless":
		return &Resource{
			kubeconfig: viper.GetString("serverless_cluster.kubeconfig"),
		}
	case "virtual_machines":
		return &Resource{
			kubeconfig: viper.GetString("virtual_machines.kubeconfig"),
		}
	}
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

func (r Resource) CreateToken(ns, serviceAccount string) (*authenticationv1.TokenRequest, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().ServiceAccounts(ns).CreateToken(context.Background(), serviceAccount, &authenticationv1.TokenRequest{}, metav1.CreateOptions{})
}
