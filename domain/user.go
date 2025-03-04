package domain

type User struct {
	id             int32
	email          string
	firstName      string
	secondName     string
	hashedPassword string
}

func NewUser(id int32, email, firstName, secondName, hashedPassword string) User {
	return User{
		id:         id,
		email:      email,
		firstName:  firstName,
		secondName: secondName,
    hashedPassword: hashedPassword,
	}
}

func (u User) ID() int32 {
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
