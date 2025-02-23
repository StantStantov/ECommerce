package stores

import (
	"Stant/ECommerce/domain"
	"database/sql"
	"fmt"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

const checkUser = `SELECT EXISTS
  (SELECT 1 FROM users 
  WHERE email = $1 
  LIMIT 1)
  ;
  `

func (s UserStore) IsExists(email string) (bool, error) {
  isExists := false
  row := s.db.QueryRow(checkUser, email)

  if err := row.Scan(&isExists); err != nil {
		return false, fmt.Errorf("UserStore IsExists: %v", err)
  }

	return isExists, nil
}

const createUser = `INSERT INTO users
  (email, first_name, second_name, password)
  VALUES
  ($1, $2, $3, $4)
  ;
  `

func (s UserStore) Create(email, fisrtName, secondName, password string) error {
	_, err := s.db.Exec(createUser, email, fisrtName, secondName, password)
	if err != nil {
		return fmt.Errorf("UserStore Create: %v", err)
	}
	return nil
}

func (s UserStore) Read(id int32) (domain.User, error) {
	return domain.User{}, nil
}
