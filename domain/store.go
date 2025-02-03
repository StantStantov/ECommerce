package domain

type ProductStore interface {
	ReadAll() ([]string, error)
}
