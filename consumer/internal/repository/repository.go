package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *Repository {
	return &Repository{
		conn: conn,
	}
}
