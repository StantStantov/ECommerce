package domain

type Product struct {
	id       int32
	price    float64
	seller   Seller
	category Category
	name     string
}

func NewProduct(id int32, name string, seller Seller, category Category, price float64) Product {
	return Product{
		id:       id,
		name:     name,
		seller:   seller,
		category: category,
		price:    price,
	}
}

func (p Product) ID() int32 {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Seller() Seller {
	return p.seller
}

func (p Product) Category() Category {
	return p.category
}

func (p Product) Price() float64 {
	return p.price
}
