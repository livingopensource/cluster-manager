package models

import "gorm.io/gorm"

type Cluster struct {
	Base
	Name       string `gorm:"varchar(255)" json:"name"`
	Kubeconfig string `gorm:"varchar(255)" json:"kubeconfig"`
}

type ClusterInterface interface {
	Save(cluster Cluster) error
	Delete(id string) error
	Update(id string, cluster Cluster) (Cluster, error)
	Find(id string) (Cluster, error)
	FindAll() ([]Cluster, error)
}

type clusterImpl struct {
	db *gorm.DB
}

func NewclusterImpl(db *gorm.DB) *clusterImpl {
	return &clusterImpl{db: db}
}

func (c clusterImpl) Save(cluster Cluster) error { return nil }

func (c clusterImpl) Delete(id string) error { return nil }

func (c clusterImpl) Update(id string, cluster Cluster) (Cluster, error) { return Cluster{}, nil }

func (c clusterImpl) Find(id string) (Cluster, error) { return Cluster{}, nil }

func (c clusterImpl) FindAll() ([]Cluster, error) { return []Cluster{}, nil }

// The codebase below this comment is used for mock test of the database

type clusterMock struct {
	db *gorm.DB
}

func NewClusterMock(db *gorm.DB) *clusterMock {
	return &clusterMock{db: db}
}

func (c clusterMock) Save(cluster Cluster) error { return nil }

func (c clusterMock) Delete(id string) error { return nil }

func (c clusterMock) Update(id string, cluster Cluster) (Cluster, error) { return Cluster{}, nil }

func (c clusterMock) Find(id string) (Cluster, error) { return Cluster{}, nil }

func (c clusterMock) FindAll() ([]Cluster, error) { return []Cluster{}, nil }
