package stores

import (
	"Stant/ECommerce/internal/domain/models"
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
  SELECT p.id, p.name, s.id, s.name, c.id, c.name, p.price
  FROM market.products p
  JOIN market.categories c ON p.category_id = c.id
  JOIN market.sellers s ON p.seller_id = s.id
  WHERE p.id = $1
  LIMIT 1
  ;
`

func (st ProductStore) Read(id string) (models.Product, error) {
	row := st.db.QueryRow(getProduct, id)

	product, err := scanProduct(row)
	if err != nil {
		return product, fmt.Errorf("stores.ProductStore.Read: [%w]", err)
	}
	return product, nil
}

const getProducts = `
  SELECT p.id, p.name, s.id, s.name, c.id, c.name, p.price
  FROM market.products p
  JOIN market.categories c ON p.category_id = c.id
  JOIN market.sellers s ON p.seller_id = s.id
  ;
`

func (st ProductStore) ReadAll() ([]models.Product, error) {
	rows, err := st.db.Query(getProducts)
	if err != nil {
		return nil, err
	}

	products := []models.Product{}
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
  SELECT p.id, p.name, s.id, s.name, c.id, c.name, p.price
  FROM market.products p
  JOIN market.categories c ON p.category_id = c.id
  JOIN market.sellers s ON p.seller_id = s.id
  WHERE
    (c.id = $1 OR NULLIF($1, '00000000-0000-0000-0000-000000000000') IS NULL)
    AND (s.id = $2 OR NULLIF($2, '00000000-0000-0000-0000-000000000000') IS NULL)
  ;
`

func (st ProductStore) ReadAllByFilter(categoryID, sellerID string) ([]models.Product, error) {
	rows, err := st.db.Query(getProductsByFilter, categoryID, sellerID)
	if err != nil {
		return nil, err
	}

	products := []models.Product{}
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

const getProductsByQuery = `
  SELECT p.id, p.name, s.id, s.name, c.id, c.name, p.price
  FROM market.products p
  JOIN market.categories c ON p.category_id = c.id
  JOIN market.sellers s ON p.seller_id = s.id
  WHERE p.name_tsv @@ to_tsquery('simple', $1)
  ;
`

func (st ProductStore) ReadAllByQuery(query string) ([]models.Product, error) {
	rows, err := st.db.Query(getProductsByQuery, query)
	if err != nil {
		return nil, err
	}

	products := []models.Product{}
	defer rows.Close()
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("stores.ProductStore.ReadAllByQuery: [%w]", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func scanProduct(row sqlRow) (models.Product, error) {
	var (
		productID    string
		name         string
		sellerID     string
		sellerName   string
		categoryID   string
		categoryName string
		price        float64
	)
	if err := row.Scan(&productID, &name, &sellerID, &sellerName, &categoryID, &categoryName, &price); err != nil {
		return models.Product{}, fmt.Errorf("stores.scanProduct: [%w]", err)
	}
	return models.NewProduct(
			productID,
			name,
			models.NewSeller(sellerID, sellerName),
			models.NewCategory(categoryID, categoryName),
			price),
		nil
}
