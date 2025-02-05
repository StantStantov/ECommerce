package main

import (
	"Stant/ECommerce/domain"
	"database/sql"
	"fmt"
)

type SQLProductStore struct {
	db *sql.DB
}

func newSQLProductStore(db *sql.DB) *SQLProductStore {
	return &SQLProductStore{db: db}
}

func (st SQLProductStore) Read(id int) (domain.Product, error) {
	q := "SELECT * FROM products WHERE product_id = $1"
	row := st.db.QueryRow(q, id)
	var productID int
	var name string
	var sellerID int
	var categoryID int
	var price float32
	if err := row.Scan(&productID, &name, &sellerID, &categoryID, &price); err != nil {
		return domain.Product{}, fmt.Errorf("SQL Read: %v", err)
	}
	return domain.NewProduct(productID, name), nil
}

func (st SQLProductStore) ReadAll() ([]domain.Product, error) {
	q := "SELECT * FROM products"
	rows, err := st.db.Query(q)
	if err != nil {
		return nil, err
	}

	products := []domain.Product{}
	defer rows.Close()
	for rows.Next() {
		var productID int
		var name string
		var sellerID int
		var categoryID int
		var price float32
		if err := rows.Scan(&productID, &name, &sellerID, &categoryID, &price); err != nil {
			return nil, err
		}
		products = append(products, domain.NewProduct(productID, name))
	}
	return products, nil
}
