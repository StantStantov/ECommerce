package stores

import (
	"Stant/ECommerce/internal/domain/models"
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
  INSERT INTO market.users
  (email, first_name, second_name, password)
  VALUES
  ($1, $2, $3, $4)
  ;
`

func (s UserStore) Create(email, fisrtName, secondName, password string) error {
	_, err := s.db.Exec(createUser, email, fisrtName, secondName, password)
	if err != nil {
		return fmt.Errorf("stores.UserStore.Create: [%w]", err)
	}
	return nil
}

const checkUser = `
  SELECT EXISTS
  (SELECT 1 FROM market.users 
  WHERE email = $1 
  LIMIT 1)
  ;
`

func (s UserStore) IsExists(email string) (bool, error) {
	var isExists bool
	row := s.db.QueryRow(checkUser, email)
	if err := row.Scan(&isExists); err != nil {
		return false, fmt.Errorf("stores.UserStore.IsExists: [%w]", err)
	}

	return isExists, nil
}

const getUserByID = `
  SELECT * FROM market.users 
  WHERE id = $1 
  LIMIT 1
  ;
`

func (s UserStore) Read(id string) (models.User, error) {
	row := s.db.QueryRow(getUserByID, id)
	user, err := scanUser(row)
	if err != nil {
		return user, fmt.Errorf("stores.UserStore.Read: [%w]", err)
	}
	return user, nil
}

const getUserByEmail = `
  SELECT * FROM market.users 
  WHERE email = $1
  LIMIT 1
  ;
`

func (s UserStore) ReadByEmail(email string) (models.User, error) {
	row := s.db.QueryRow(getUserByEmail, email)
	user, err := scanUser(row)
	if err != nil {
		return user, fmt.Errorf("stores.UserStore.ReadByEmail: [%w]", err)
	}
	return user, nil
}

func scanUser(row sqlRow) (models.User, error) {
	var (
		id             string
		email          string
		firstName      string
		secondName     string
		hashedPassword string
	)
	if err := row.Scan(&id, &email, &firstName, &secondName, &hashedPassword); err != nil {
		return models.User{}, fmt.Errorf("stores.scanUser: [%w]", err)
	}
	return models.NewUser(id, email, firstName, secondName, hashedPassword), nil
}
