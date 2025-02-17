package stores

import "Stant/ECommerce/domain"

type sqlRow interface {
	Scan(dest ...any) error
}

func scanProduct(row sqlRow) (domain.Product, error) {
	var productID int32
	var name string
	var sellerID int32
	var sellerName string
	var categoryID int32
	var categoryName string
	var price float64
	if err := row.Scan(&productID, &name, &sellerID, &sellerName, &categoryID, &categoryName, &price); err != nil {
		return domain.Product{}, err
	}
	return domain.NewProduct(productID, name, domain.NewSeller(sellerID, sellerName), domain.NewCategory(categoryID, categoryName), price), nil
}

func scanCategory(row sqlRow) (domain.Category, error) {
	var id int32
	var name string
	if err := row.Scan(&id, &name); err != nil {
		return domain.Category{}, err
	}
	return domain.NewCategory(id, name), nil
}

func scanSeller(row sqlRow) (domain.Seller, error) {
	var id int32
	var name string
	if err := row.Scan(&id, &name); err != nil {
		return domain.Seller{}, err
	}
	return domain.NewSeller(id, name), nil
}
