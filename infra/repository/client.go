package repository

import (
	"database/sql"
	"os"
)

type Client struct {
	db *sql.DB
}

func NewClient() (*Client, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}
