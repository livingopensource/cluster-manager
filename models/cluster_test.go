package models

import (
	"testing"

	"gorm.io/gorm"
)

func ClusterTest(t *testing.T) {
	db := &gorm.DB{}
	cluster := NewClusterMock(db)

	err := cluster.Save(Cluster{Name: "db-cluster", Kubeconfig: "This is a configuration file"})
	if err != nil {
		t.Errorf("Expected a nil error, but got %s", err.Error())
	}
}
