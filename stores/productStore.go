package stores

import (
	"Stant/ECommerce/domain"
	"database/sql"
	"fmt"
)

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

func (st ProductStore) Read(id int) (domain.Product, error) {
	q := "SELECT * FROM products WHERE product_id = $1"
	row := st.db.QueryRow(q, id)
	var productID int
	var name string
	var sellerID int
	var categoryID int
	var price float32
	if err := row.Scan(&productID, &name, &sellerID, &categoryID, &price); err != nil {
		return domain.Product{}, fmt.Errorf(" Read: %v", err)
	}
	return domain.NewProduct(productID, name, sellerID, categoryID, price), nil
}

func (st ProductStore) ReadAll() ([]domain.Product, error) {
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
		products = append(products, domain.NewProduct(productID, name, sellerID, categoryID, price))
	}
	return products, nil
}

func (st ProductStore) ReadAllByFilter(categoryID int) ([]domain.Product, error) {
	q := "SELECT * FROM products WHERE category_id = $1"
	rows, err := st.db.Query(q, categoryID)
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
		products = append(products, domain.NewProduct(productID, name, sellerID, categoryID, price))
	}
	return products, nil
}
