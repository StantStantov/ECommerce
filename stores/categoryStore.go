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

func (st CategoryStore) Read(categoryID int) (domain.Category, error) {
	q := `
  SELECT * 
  FROM categories 
  WHERE category_id = $1
  ;
  `
	row := st.db.QueryRow(q, categoryID)
	category, err := scanCategory(row)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (st CategoryStore) ReadAll() ([]domain.Category, error) {
	q := `
  SELECT * 
  FROM categories
  ;
  `
	rows, err := st.db.Query(q)
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
