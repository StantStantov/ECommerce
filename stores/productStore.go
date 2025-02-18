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

const getProduct = `
  SELECT p.product_id, p.product_name, s.seller_id, s.seller_name, c.category_id, c.category_name, p.product_price
  FROM products p
  JOIN categories c ON p.category_id = c.category_id
  JOIN sellers s ON p.seller_id = s.seller_id
  WHERE p.product_id = $1
  LIMIT 1
  ;
  `

func (st ProductStore) Read(id int) (domain.Product, error) {
	row := st.db.QueryRow(getProduct, id)

	product, err := scanProduct(row)
	if err != nil {
		return domain.Product{}, fmt.Errorf("ProductStore Read: %v", err)
	}
	return product, nil
}

const getProducts = `
  SELECT p.product_id, p.product_name, s.seller_id, s.seller_name, c.category_id, c.category_name, p.product_price
  FROM products p
  JOIN categories c ON p.category_id = c.category_id
  JOIN sellers s ON p.seller_id = s.seller_id
  ;
  `

func (st ProductStore) ReadAll() ([]domain.Product, error) {
	rows, err := st.db.Query(getProducts)
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

const getProductsByFilter = `
  SELECT p.product_id, p.product_name, s.seller_id, s.seller_name, c.category_id, c.category_name, p.product_price
  FROM products p
  JOIN categories c ON p.category_id = c.category_id
  JOIN sellers s ON p.seller_id = s.seller_id
  WHERE
    (c.category_id = $1 OR NULLIF($1, 0) IS NULL)
    AND (s.seller_id = $2 OR NULLIF($2, 0) IS NULL)
  ;
  `

func (st ProductStore) ReadAllByFilter(categoryID int, sellerID int) ([]domain.Product, error) {
	rows, err := st.db.Query(getProductsByFilter, categoryID, sellerID)
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
