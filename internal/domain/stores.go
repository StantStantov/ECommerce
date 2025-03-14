package domain

import (
	"log"
	"time"
)

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
	ReadByEmail(email string) (User, error)
}

type SessionStore interface {
	Create(sessionToken, csrfToken string, expireOn time.Time) error
	Read(sessionToken string) (Session, error)
	Delete(sessionToken string) error
	DeleteAllExpired() error
	StartCleanup(logger log.Logger, interval time.Duration) (chan<- struct{}, <-chan struct{})
	StopCleanup(quit chan<- struct{}, done <-chan struct{})
}
