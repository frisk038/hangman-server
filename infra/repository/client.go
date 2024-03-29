package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/frisk038/hangman-server/business"
	"github.com/frisk038/hangman-server/business/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client struct {
	db *pgxpool.Pool
}

const selectYesterdayNum = "SELECT NUM FROM SECRET ORDER BY SECRETID DESC LIMIT 1;"
const insertTodaySecret = "INSERT INTO SECRET (NUM, VALUE) VALUES($1, $2) ON CONFLICT (NUM) DO NOTHING;"
const selectTodaySecret = "SELECT NUM, VALUE FROM SECRET ORDER BY SECRETID DESC LIMIT 1;"
const selectYesterdaySecret = "SELECT PICKEDDT, NUM FROM SECRET ORDER BY SECRETID DESC LIMIT 1;"
const insertUserScore = "INSERT INTO USERSCORE (USERID, SECRETNUM, SCORE, NAME, USERAGENT) VALUES ($1, $2, $3, NULLIF($4, ''), $5);"
const updateUserName = "UPDATE userscore SET name = COALESCE($1, name) WHERE userid = $2 AND secretnum = $3 AND name IS NULL RETURNING userid;"
const selectTopPlayer = "SELECT name, SUM(score) AS highScore FROM userscore WHERE playdt >= DATE_TRUNC('week',NOW()) AND name is not NULL GROUP BY name ORDER BY highScore DESC LIMIT 5;"
const selectWeekTopPlayer = "SELECT name, SUM(score) AS highScore FROM userscore WHERE secretnum >= $1 - 7 and secretnum < $1 AND name is not NULL GROUP BY name ORDER BY highScore DESC LIMIT 1;"

// const selectTopPlayer = "select name, score from userscore where SECRETNUM = $1 AND NAME is not null  ORDER BY score DESC LIMIT 5;"

// const updateUserName = "UPDATE USERSCORE SET NAME = $1 WHERE USERID = $2 AND SECRETNUM = $3 AND name IS NULL RETURNING USERID;"
// const selectTopPlayer = "SELECT DISTINCT ON (score) name, score FROM userscore where SECRETNUM = $1 and name is NOT NULL ORDER BY score LIMIT 3;"

func NewClient() (*Client, error) {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
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

func (c *Client) SelectYesterdaySecret(ctx context.Context) (entity.Secret, error) {
	var secret entity.Secret
	err := c.db.QueryRow(ctx, selectYesterdaySecret).Scan(&secret.PickedDt, &secret.Number)
	return secret, err
}

func (c *Client) InsertUserScore(ctx context.Context, score entity.Score) error {
	row, _ := c.db.Query(ctx, insertUserScore, score.UserID, score.SecretNum, score.Score, score.UserName, score.UserAgent)
	defer row.Close()
	return row.Err()
}

func (c *Client) UpdateUserName(ctx context.Context, score entity.Score) error {
	row, err := c.db.Query(ctx, updateUserName, score.UserName, score.UserID, score.SecretNum)
	if err != nil {
		return err
	}
	defer row.Close()

	nbRow := 0
	for row.Next() {
		nbRow++
	}

	switch {
	case nbRow > 1:
		return fmt.Errorf("insert user went wrong : (%s,%d)", score.UserID.String(), score.SecretNum)
	case nbRow == 0:
		return business.AlreadyExistUserErr
	default:
		return row.Err()
	}
}

func (c *Client) SelectTopPlayer(ctx context.Context) ([]entity.Score, error) {
	row, err := c.db.Query(ctx, selectTopPlayer)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	topN := []entity.Score{}
	for row.Next() {
		var top entity.Score
		err := row.Scan(&top.UserName, &top.Score)
		if err != nil {
			return nil, err
		}
		topN = append(topN, top)
	}
	return topN, row.Err()
}

func (c *Client) SelectWeeklyTopPlayer(ctx context.Context, secretNum int) (entity.Score, error) {
	var top entity.Score
	err := c.db.QueryRow(ctx, selectWeekTopPlayer, secretNum).Scan(&top.UserName, &top.Score)
	if err != nil {
		return entity.Score{}, err
	}

	return top, nil
}
