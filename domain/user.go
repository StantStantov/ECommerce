package domain

type User struct {
	id         int32
	email      string
	firstName  string
	secondName string
}

func newUser(id int32, email, firstName, secondName string) User {
	return User{
		id:         id,
		email:      email,
		firstName:  firstName,
		secondName: secondName,
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
