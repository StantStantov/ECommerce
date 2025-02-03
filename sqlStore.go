package main

import (
	"database/sql"
)

type SQLProductStore struct {
	db *sql.DB
}

func newSQLProductStore(db *sql.DB) *SQLProductStore {
	return &SQLProductStore{db: db}
}

func (st SQLProductStore) ReadAll() ([]string, error) {
	q := "SELECT name FROM laptops"
	rows, err := st.db.Query(q)
	if err != nil {
		return nil, err
	}

	products := []string{}
	defer rows.Close()
	for rows.Next() {
		var product string
		if err := rows.Scan(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
