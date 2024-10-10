package clusters

import "net/url"

type ClusterResource struct {
	ID        string  `json:"id,omitempty"`
	Namespace string  `json:"namespace,omitempty"`
	Compute   Compute `json:"compute,omitempty"`
	Account   Account `json:"account,omitempty"`
	HTTP      HTTP    `json:"http,omitempty"`
}

type Compute struct {
	Name      string      `json:"name,omitempty"`
	CPU       float64     `json:"vcpu,omitempty"`
	RAM       string      `json:"ram,omitempty"`
	Storage   string      `json:"storage,omitempty"`
	Instances float64     `json:"instances,omitempty"`
	State     string      `json:"state,omitempty"`
	SSHKey    string      `json:"ssh_key,omitempty"`
	URL       string      `json:"url,omitempty"`
	Container []Container `json:"containers,omitempty"`
}

type Account struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type HTTP struct {
	QueryParams url.Values `json:"params,omitempty"`
}

type Container struct {
	Image string `json:"image,omitempty"`
	Port  []Port `json:"ports,omitempty"`
	Env   []Env  `json:"env,omitempty"`
}

type Port struct {
	ContainerPort int32 `json:"containerPort,omitempty"`
}

type Env struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
