package domain

type Product struct {
	name string
}

func NewProduct(name string) Product {
	return Product{name: name}
}

func (p Product) GetName() string {
	return p.name
}
