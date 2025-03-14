package domain

type Seller struct {
	id   int32
	name string
}

func NewSeller(id int32, name string) Seller {
	return Seller{
		id:   id,
		name: name,
	}
}

func (s Seller) ID() int32 {
	return s.id
}

func (s Seller) Name() string {
	return s.name
}
