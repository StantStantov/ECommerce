package domain

type Seller struct {
	id   int
	name string
}

func NewSeller(id int, name string) Seller {
	return Seller{
		id:   id,
		name: name,
	}
}

func (s Seller) ID() int {
	return s.id
}

func (s Seller) Name() string {
	return s.name
}
