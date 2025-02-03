package domain

type ProductStore interface {
	ReadAll() ([]Product, error)
}
