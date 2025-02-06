package stores

import (
	"Stant/ECommerce/domain"
	"database/sql"
	"fmt"
)

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{db: db}
}

func (st CategoryStore) Read(categoryID int) (domain.Category, error) {
	q := "SELECT * FROM categories WHERE category_id = $1"
	row := st.db.QueryRow(q, categoryID)
	var id int
	var name string
	if err := row.Scan(&id, &name); err != nil {
		return domain.Category{}, fmt.Errorf("SQL Read: %v", err)
	}
	return domain.NewCategory(id, name), nil
}

func (st CategoryStore) ReadAll() ([]domain.Category, error) {
	q := "SELECT * FROM categories"
	rows, err := st.db.Query(q)
	if err != nil {
		return nil, err
	}

	categories := []domain.Category{}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("SQL ReadAll: %v", err)
		}
		categories = append(categories, domain.NewCategory(id, name))
	}
	return categories, nil
}
