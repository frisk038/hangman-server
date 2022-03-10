package repository

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

type Client struct {
	db *pgx.Conn
}

func NewClient() (*Client, error) {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}
