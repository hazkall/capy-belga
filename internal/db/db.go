package db

import (
	"database/sql"
)

func (p *Postgres) Execute(query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

func (p *Postgres) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.Query(query, args...)
}

func (p *Postgres) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRow(query, args...)
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}
