package stores

import "Stant/ECommerce/domain"

type sqlRow interface {
	Scan(dest ...any) error
}

func scanProduct(row sqlRow) (domain.Product, error) {
	var productID int
	var name string
	var sellerID int
	var categoryID int
	var price float32
	if err := row.Scan(&productID, &name, &sellerID, &categoryID, &price); err != nil {
		return domain.Product{}, err
	}
	return domain.NewProduct(productID, name, sellerID, categoryID, price), nil
}

func scanCategory(row sqlRow) (domain.Category, error) {
	var id int
	var name string
	if err := row.Scan(&id, &name); err != nil {
		return domain.Category{}, err
	}
	return domain.NewCategory(id, name), nil
}
