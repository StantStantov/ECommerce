package domain

type Product struct {
	id         int
	sellerId   int
	categoryId int
	name       string
	price      float32
}

func NewProduct(id int, name string, sellerId int, categoryId int, price float32) Product {
	return Product{
		id:         id,
		name:       name,
		sellerId:   sellerId,
		categoryId: categoryId,
		price:      price,
	}
}

func (p Product) ID() int {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) SellerID() int {
	return p.sellerId
}

func (p Product) CategoryID() int {
	return p.categoryId
}

func (p Product) Price() int {
	return p.categoryId
}
