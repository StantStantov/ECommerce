package models

type User struct {
	id             string
	email          string
	firstName      string
	secondName     string
	hashedPassword string
}

func NewUser(id, email, firstName, secondName, hashedPassword string) User {
	return User{
		id:             id,
		email:          email,
		firstName:      firstName,
		secondName:     secondName,
		hashedPassword: hashedPassword,
	}
}

func (u User) ID() string {
	return u.id
}

func (u User) Email() string {
	return u.email
}

func (u User) FirstName() string {
	return u.firstName
}

func (u User) SecondName() string {
	return u.secondName
}

func (u User) HashedPassword() string {
	return u.hashedPassword
}
