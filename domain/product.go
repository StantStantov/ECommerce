package domain

type Product struct {
	id       int
	price    float64
	seller   string
	category string
	name     string
}

func NewProduct(id int, name string, seller string, category string, price float64) Product {
	return Product{
		id:       id,
		name:     name,
		seller:   seller,
		category: category,
		price:    price,
	}
}

func (p Product) ID() int {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Seller() string {
	return p.seller
}

func (p Product) Category() string {
	return p.category
}

func (p Product) Price() float64 {
	return p.price
}
