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

type SellerStore interface {
	Read(id int) (Seller, error)
	ReadAll() ([]Seller, error)
}

type UserStore interface {
	IsExists(email string) (bool, error)
	Create(email, fisrtName, secondName, password string) error
	Read(id int32) (User, error)
}
