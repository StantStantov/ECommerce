package main

import (
	"Stant/ECommerce/domain"
)

type InMemoryStore struct {
	db []domain.Product
}

func newInMemoryStore(db []domain.Product) *InMemoryStore {
	return &InMemoryStore{db: db}
}

func (st InMemoryStore) Read(id int) (domain.Product, error) {
	return st.db[id], nil
}

func (st InMemoryStore) ReadAll() ([]domain.Product, error) {
	return st.db, nil
}

func (st InMemoryStore) ReadAllByFilter(categoryID int) ([]domain.Product, error) {
	products := []domain.Product{}
	for _, product := range st.db {
		if product.CategoryID() == categoryID {
			products = append(products, product)
		}
	}
	return products, nil
}
