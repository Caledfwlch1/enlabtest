package postgres

import (
	"database/sql"
	"fmt"

	"github.com/caledfwlch1/enlabtest/db"
	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func NewDatabase(host, user, password, dbName, options string) (db.Database, error) {
	//"postgres://pqgotest:password@localhost/pqgotest?sslmode=disable"
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?%s", user, password, host, dbName, options)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opnen database %s", err)
	}

	return &postgres{
		db: db,
	}, nil
}
