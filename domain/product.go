package domain

type Product struct {
	id   int
	name string
}

func NewProduct(id int, name string) Product {
	return Product{id: id, name: name}
}

func (p Product) GetID() int {
	return p.id
}

func (p Product) GetName() string {
	return p.name
}
