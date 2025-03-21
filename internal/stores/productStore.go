package stores

import (
	"Stant/ECommerce/internal/domain"
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
		return product, fmt.Errorf("stores.ProductStore.Read: [%w]", err)
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
			return nil, fmt.Errorf("stores.ProductStore.ReadAll: [%w]", err)
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
			return nil, fmt.Errorf("stores.ProductStore.ReadAllByFilter: [%w]", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func scanProduct(row sqlRow) (domain.Product, error) {
	var (
		productID    int32
		name         string
		sellerID     int32
		sellerName   string
		categoryID   int32
		categoryName string
		price        float64
	)
	if err := row.Scan(&productID, &name, &sellerID, &sellerName, &categoryID, &categoryName, &price); err != nil {
		return domain.Product{}, fmt.Errorf("stores.scanProduct: [%w]", err)
	}
	return domain.NewProduct(
			productID,
			name,
			domain.NewSeller(sellerID, sellerName),
			domain.NewCategory(categoryID, categoryName),
			price),
		nil
}
