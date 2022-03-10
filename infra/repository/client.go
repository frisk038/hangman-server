package repository

import (
	"context"
	"os"

	"github.com/frisk038/hangman-server/business/entity"
	"github.com/jackc/pgx/v4"
)

type Client struct {
	db *pgx.Conn
}

const selectYesterdayNum = "SELECT NUM FROM SECRET ORDER BY SECRETID DESC LIMIT 1;"
const insertTodaySecret = "INSERT INTO SECRET (NUM, VALUE) VALUES ($1, $2);"
const selectTodaySecret = "SELECT NUM, VALUE FROM SECRET ORDER BY SECRETID DESC LIMIT 1;"

func NewClient() (*Client, error) {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) GetYesterdayNumber(ctx context.Context) (int, error) {
	var num int
	err := c.db.QueryRow(ctx, selectYesterdayNum).Scan(&num)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return num, err
}

func (c *Client) InsertTodaySecret(ctx context.Context, secret entity.Secret) error {
	row, _ := c.db.Query(ctx, insertTodaySecret, secret.Number, secret.SecretWord)
	defer row.Close()
	return row.Err()
}

func (c *Client) SelectTodaySecret(ctx context.Context) (entity.Secret, error) {
	var secret entity.Secret
	err := c.db.QueryRow(ctx, selectTodaySecret).Scan(&secret.Number, &secret.SecretWord)
	return secret, err
}
