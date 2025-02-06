package stores

import (
	"database/sql"
	"fmt"
)

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{db: db}
}

func (st CategoryStore) Read(id int) (string, error) {
	q := "SELECT * FROM categories WHERE category_id = $1"
	row := st.db.QueryRow(q, id)
	var category string
	if err := row.Scan(category); err != nil {
		return "", fmt.Errorf("SQL Read: %v", err)
	}
	return category, nil
}

func (st CategoryStore) ReadAll() ([]string, error) {
	q := "SELECT * FROM products"
	rows, err := st.db.Query(q)
	if err != nil {
		return nil, err
	}

	categories := []string{}
	defer rows.Close()
	for rows.Next() {
		var category string
		if err := rows.Scan(category); err != nil {
			return nil, fmt.Errorf("SQL ReadAll: %v", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
