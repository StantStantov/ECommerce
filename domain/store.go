package domain

type ProductStore interface {
	Read(id int) (Product, error)
	ReadAll() ([]Product, error)
}
