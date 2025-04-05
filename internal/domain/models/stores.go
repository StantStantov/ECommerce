package models

import (
	"log"
	"time"
)

type ProductStore interface {
	Read(id string) (Product, error)
	ReadAll() ([]Product, error)
	ReadAllByFilter(categoryID string, sellerID string) ([]Product, error)
}

type CategoryStore interface {
	Read(id string) (Category, error)
	ReadAll() ([]Category, error)
}

type SellerStore interface {
	Read(id string) (Seller, error)
	ReadAll() ([]Seller, error)
}

type UserStore interface {
	IsExists(email string) (bool, error)
	Create(email, fisrtName, secondName, password string) error
	Read(id string) (User, error)
	ReadByEmail(email string) (User, error)
}

type SessionStore interface {
	Create(userID, sessionToken, csrfToken string) error
	Read(sessionToken string) (Session, error)
	Delete(sessionToken string) error
	DeleteAllExpired() error
	StartCleanup(logger log.Logger, interval time.Duration) (chan<- struct{}, <-chan struct{})
	StopCleanup(quit chan<- struct{}, done <-chan struct{})
}
