package dao

import (
	"database/sql"
	"fmt"

	"go_code/simplek8s/core/application/repository"
	"go_code/simplek8s/core/entity"
)

type clusterDao struct {
	DB *sql.DB
}

func NewClusterDao(db *sql.DB) repository.ClusterRepo {
	return &clusterDao{DB: db}
}

func (dao *clusterDao) Create(cluster entity.Cluster) (int64, error) {
	stmt, err := dao.DB.Prepare("INSERT INTO clusters(config) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(cluster.Config)
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

func (dao *clusterDao) GetByID(id int) (entity.Cluster, error) {
	var cluster entity.Cluster
	err := dao.DB.QueryRow("SELECT id, config FROM clusters WHERE id = ?", id).Scan(&cluster.ID, &cluster.Config)
	if err != nil {
		if err == sql.ErrNoRows {
			return cluster, fmt.Errorf("no cluster found with id %d", id)
		}
		return cluster, fmt.Errorf("failed to query row: %v", err)
	}

	return cluster, nil
}

func (dao *clusterDao) GetAll() ([]entity.Cluster, error) {
	rows, err := dao.DB.Query("SELECT id, config FROM clusters")
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %v", err)
	}
	defer rows.Close()

	var clusters []entity.Cluster
	for rows.Next() {
		var cluster entity.Cluster
		err := rows.Scan(&cluster.ID, &cluster.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		clusters = append(clusters, cluster)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return clusters, nil
}
