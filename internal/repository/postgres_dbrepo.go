package repository

import (
	"database/sql"
	"time"
)

const dbTimeOut = time.Second * 3

type PostgresDBRepo struct {
	DB *sql.DB
}

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}
