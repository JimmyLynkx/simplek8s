package repository

import "go_code/simplek8s/core/entity"

type ClusterRepo interface {
	Create(cluster entity.Cluster) (int64, error)
	GetByID(id int) (entity.Cluster, error)
	GetAll() ([]entity.Cluster, error)
}
