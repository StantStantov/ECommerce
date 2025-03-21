package stores

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func NewDBConn() (*sql.DB, error) {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok || len(user) == 0 {
		return nil, errors.New("stores.NewDBConn: [Env POSTGRES_USER doesn't exist or empty]")
	}
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok || len(password) == 0 {
		return nil, errors.New("stores.NewDBConn: [Env POSTGRES_PASSWORD doesn't exist or empty]")
	}
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok || len(host) == 0 {
		return nil, errors.New("stores.NewDBConn: [Env POSTGRES_HOST doesn't exist or empty]")
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok || len(port) == 0 {
		return nil, errors.New("stores.NewDBConn: [Env POSTGRES_PORT doesn't exist or empty]")
	}
	dbname, ok := os.LookupEnv("POSTGRES_DB")
	if !ok || len(dbname) == 0 {
		return nil, errors.New("stores.NewDBConn: [Env POSTGRES_DB doesn't exist or empty]")
	}

	pgConn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		user, password, host, port, dbname,
	)
	db, err := sql.Open("pgx", pgConn)
	if err != nil {
		return nil, fmt.Errorf("stores.NewDBConn: [%w]", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("stores.NewDBConn: [%w]", err)
	}
	return db, nil
}

type sqlRow interface {
	Scan(dest ...any) error
}
