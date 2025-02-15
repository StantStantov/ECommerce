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
	q := `
  SELECT p.product_id, p.product_name, s.seller_name, c.category_name, p.product_price
  FROM products p
  JOIN categories c ON p.category_id = c.category_id
  JOIN sellers s ON p.seller_id = s.seller_id
  WHERE p.product_id = $1
  ;
  `
	row := st.db.QueryRow(q, id)

	product, err := scanProduct(row)
	if err != nil {
		return domain.Product{}, fmt.Errorf("ProductStore Read: %v", err)
	}
	return product, nil
}

func (st ProductStore) ReadAll() ([]domain.Product, error) {
	q := `
  SELECT p.product_id, p.product_name, s.seller_name, c.category_name, p.product_price
  FROM products p
  JOIN categories c ON p.category_id = c.category_id
  JOIN sellers s ON p.seller_id = s.seller_id
  ;
  `
	rows, err := st.db.Query(q)
	if err != nil {
		return nil, err
	}

	products := []domain.Product{}
	defer rows.Close()
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("ProductStore ReadAll: %v", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (st ProductStore) ReadAllByFilter(categoryID int) ([]domain.Product, error) {
	q := `
  SELECT p.product_id, p.product_name, s.seller_name, c.category_name, p.product_price
  FROM products p
  JOIN categories c ON p.category_id = c.category_id
  JOIN sellers s ON p.seller_id = s.seller_id
  WHERE c.category_id = $1
  ;
  `
	rows, err := st.db.Query(q, categoryID)
	if err != nil {
		return nil, err
	}

	products := []domain.Product{}
	defer rows.Close()
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("ProductStore ReadAllByFilter: %v", err)
		}
		products = append(products, product)
	}
	return products, nil
}
