package domain

type ProductStore interface {
	Read(id int) (Product, error)
	ReadAll() ([]Product, error)
	ReadAllByFilter(categoryID int) ([]Product, error)
}

type CategoryStore interface {
	Read(id int) (string, error)
	ReadAll() ([]string, error)
}
