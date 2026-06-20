package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Open opens and verifies a MySQL connection.
func Open(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}
