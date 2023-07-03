package config

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func OpenDatabase() error {
	var err error
	connStr := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return nil
}

func GetDBConn() *sqlx.DB {
	return DB
}

func CloseDatabase() error {
	return DB.Close()
}
