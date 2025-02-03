package main

type InMemoryStore struct {
	db []string
}

func newInMemoryStore(db []string) *InMemoryStore {
	return &InMemoryStore{db: db}
}

func (st InMemoryStore) ReadAll() ([]string, error) {
	return st.db, nil
}
