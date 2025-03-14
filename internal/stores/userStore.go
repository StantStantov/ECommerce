package stores

import (
	"Stant/ECommerce/internal/domain"
	"database/sql"
	"fmt"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

const createUser = `
  INSERT INTO users
  (email, first_name, second_name, password)
  VALUES
  ($1, $2, $3, $4)
  ;
`

func (s UserStore) Create(email, fisrtName, secondName, password string) error {
	_, err := s.db.Exec(createUser, email, fisrtName, secondName, password)
	if err != nil {
		return fmt.Errorf("UserStore.Create: [%w]", err)
	}
	return nil
}

const checkUser = `
  SELECT EXISTS
  (SELECT 1 FROM users 
  WHERE email = $1 
  LIMIT 1)
  ;
`

func (s UserStore) IsExists(email string) (bool, error) {
	isExists := false
	row := s.db.QueryRow(checkUser, email)

	if err := row.Scan(&isExists); err != nil {
		return false, fmt.Errorf("UserStore.IsExists: [%w]", err)
	}

	return isExists, nil
}

const getUserByID = `
  SELECT * FROM users 
  WHERE id = $1 
  LIMIT 1
  ;
`

func (s UserStore) Read(id int32) (domain.User, error) {
	row := s.db.QueryRow(getUserByID, id)
	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("UserStore.Read: [%w]", err)
	}
	return user, nil
}

const getUserByEmail = `
  SELECT * FROM users 
  WHERE email = $1
  LIMIT 1
  ;
`

func (s UserStore) ReadByEmail(email string) (domain.User, error) {
	row := s.db.QueryRow(getUserByEmail, email)
	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("UserStore.ReadByEmail: [%w]", err)
	}
	return user, nil
}

func scanUser(row sqlRow) (domain.User, error) {
	var (
		id             int32
		email          string
		firstName      string
		secondName     string
		hashedPassword string
	)
	if err := row.Scan(&id, &email, &firstName, &secondName, &hashedPassword); err != nil {
		return domain.User{}, fmt.Errorf("scanUser: [%w]", err)
	}
	return domain.NewUser(id, email, firstName, secondName, hashedPassword), nil
}
