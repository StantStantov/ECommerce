package stores

import "Stant/ECommerce/domain"

type sqlRow interface {
	Scan(dest ...any) error
}

func scanProduct(row sqlRow) (domain.Product, error) {
	var productID int
	var name string
	var seller string
	var category string
	var price float64
	if err := row.Scan(&productID, &name, &seller, &category, &price); err != nil {
		return domain.Product{}, err
	}
	return domain.NewProduct(productID, name, seller, category, price), nil
}

func scanCategory(row sqlRow) (domain.Category, error) {
	var id int
	var name string
	if err := row.Scan(&id, &name); err != nil {
		return domain.Category{}, err
	}
	return domain.NewCategory(id, name), nil
}
