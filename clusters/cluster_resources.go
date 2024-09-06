package clusters

type ClusterResource struct {
	ID        string  `json:"id,omitempty"`
	Namespace string  `json:"namespace,omitempty"`
	Compute   Compute `json:"compute"`
	Account   Account `json:"account"`
}

type Compute struct {
	Name      string  `json:"name,omitempty"`
	CPU       float64 `json:"vcpu,omitempty"`
	RAM       float64 `json:"ram,omitempty"`
	Storage   float64 `json:"storage,omitempty"`
	Instances float64 `json:"instances,omitempty"`
}

type Account struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}
