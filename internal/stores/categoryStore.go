package stores

import (
	"Stant/ECommerce/internal/domain"
	"database/sql"
	"fmt"
)

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{db: db}
}

const getCategory = `
  SELECT * 
  FROM market.categories 
  WHERE id = $1 
  LIMIT 1
  ;
`

func (st CategoryStore) Read(categoryID string) (domain.Category, error) {
	row := st.db.QueryRow(getCategory, categoryID)
	category, err := scanCategory(row)
	if err != nil {
		return category, fmt.Errorf("stores.CategoryStore.Read: [%w]", err)
	}
	return category, nil
}

const getCategories = `
  SELECT * 
  FROM market.categories
  ;
`

func (st CategoryStore) ReadAll() ([]domain.Category, error) {
	rows, err := st.db.Query(getCategories)
	if err != nil {
		return nil, err
	}

	categories := []domain.Category{}
	defer rows.Close()
	for rows.Next() {
		category, err := scanCategory(rows)
		if err != nil {
			return nil, fmt.Errorf("stores.CategoryStore.ReadAll: [%w]", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func scanCategory(row sqlRow) (domain.Category, error) {
	var (
		id   string
		name string
	)
	if err := row.Scan(&id, &name); err != nil {
		return domain.Category{}, fmt.Errorf("stores.scanCategory: [%w]", err)
	}
	return domain.NewCategory(id, name), nil
}
