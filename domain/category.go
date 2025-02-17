package domain

type Category struct {
	id   int32
	name string
}

func NewCategory(id int32, name string) Category {
	return Category{id: id, name: name}
}

func (c Category) ID() int32 {
  return c.id
}

func (c Category) Name() string {
  return c.name
}
