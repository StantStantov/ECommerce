package stores

import (
	"database/sql"
	"fmt"
	"os"
)

func NewDBConn() (*sql.DB, error) {
	pgConn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("pgx", pgConn)
	return db, err
}

type sqlRow interface {
	Scan(dest ...any) error
}
