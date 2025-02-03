package main

import (
	"Stant/ECommerce/domain"
	"database/sql"
)

type SQLProductStore struct {
	db *sql.DB
}

func newSQLProductStore(db *sql.DB) *SQLProductStore {
	return &SQLProductStore{db: db}
}

func (st SQLProductStore) ReadAll() ([]domain.Product, error) {
	q := "SELECT name FROM laptops"
	rows, err := st.db.Query(q)
	if err != nil {
		return nil, err
	}

	products := []domain.Product{}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		products = append(products, domain.NewProduct(name))
	}
	return products, nil
}
