package domain

type ProductStore interface {
	Read(id int) (Product, error)
	ReadAll() ([]Product, error)
	ReadAllByFilter(categoryID int, sellerID int) ([]Product, error)
}

type CategoryStore interface {
	Read(id int) (Category, error)
	ReadAll() ([]Category, error)
}
