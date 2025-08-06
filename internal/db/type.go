package db

import "database/sql"

type Postgres struct {
	DB *sql.DB
}
