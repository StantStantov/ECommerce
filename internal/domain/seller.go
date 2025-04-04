package domain

type Seller struct {
	id   string
	name string
}

func NewSeller(id, name string) Seller {
	return Seller{
		id:   id,
		name: name,
	}
}

func (s Seller) ID() string {
	return s.id
}

func (s Seller) Name() string {
	return s.name
}
