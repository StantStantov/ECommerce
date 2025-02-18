package stores

import (
	"Stant/ECommerce/domain"
	"database/sql"
)

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{db: db}
}

const getCategory = `
  SELECT * 
  FROM categories 
  WHERE category_id = $1 LIMIT 1
  ;
  `

func (st CategoryStore) Read(categoryID int) (domain.Category, error) {
	row := st.db.QueryRow(getCategory, categoryID)
	category, err := scanCategory(row)
	if err != nil {
		return category, err
	}
	return category, nil
}

const getCategories = `
  SELECT * 
  FROM categories
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
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
