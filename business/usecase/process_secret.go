package usecase

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/frisk038/hangman-server/business/entity"
)

const maxNbWord = 159829

type repository interface {
	GetYesterdayNumber(ctx context.Context) (int, error)
	InsertTodaySecret(ctx context.Context, secret entity.Secret) error
	SelectTodaySecret(ctx context.Context) (entity.Secret, error)
	InsertUserScore(ctx context.Context, score entity.Score) error
	UpdateUserName(ctx context.Context, score entity.Score) error
	SelectTopPlayer(ctx context.Context) ([]entity.Score, error)
	SelectYesterdaySecret(ctx context.Context) (entity.Secret, error)
}

type ProcessSecret struct {
	repo repository
}

func NewProcessSecret(repo repository) ProcessSecret {
	rand.Seed(time.Now().UnixNano())
	return ProcessSecret{repo: repo}
}

func (ps ProcessSecret) InsertSecretTask() {
	fmt.Println("InsertSecretTask runing...")
	secret, err := ps.generateDailySecret()
	if err != nil {
		log.Print(err)
	}

	nb, err := ps.repo.GetYesterdayNumber(context.Background())
	if err != nil {
		log.Print(err)
	}
	secret.Number = nb + 1

	err = ps.repo.InsertTodaySecret(context.Background(), secret)
	if err != nil {
		log.Print(err)
	}
}

func (ps ProcessSecret) GetSecret(ctx context.Context) (entity.Secret, error) {
	previousSecret, err := ps.repo.SelectYesterdaySecret(ctx)
	if err != nil {
		return entity.Secret{}, err
	}

	if previousSecret.PickedDt.Day() <= time.Now().Add(-24*time.Hour).Day() {
		secret, err := ps.generateDailySecret()
		if err != nil {
			log.Print(err)
		}
		secret.Number = previousSecret.Number + 1

		err = ps.repo.InsertTodaySecret(context.Background(), secret)
		if err != nil {
			log.Print(err)
		}
	}
	return ps.repo.SelectTodaySecret(ctx)
}

func (ps ProcessSecret) generateDailySecret() (entity.Secret, error) {
	f, err := os.Open("./business/usecase/dico.txt")
	if err != nil {
		dir, _ := os.Getwd()
		return entity.Secret{}, fmt.Errorf(" %w (%s)", err, dir)
	}
	defer f.Close()

	var line int
	scanner := bufio.NewScanner(f)
	wantedLine := rand.Intn(maxNbWord)

	for scanner.Scan() {
		if line == wantedLine {
			return entity.Secret{
				SecretWord: scanner.Text(),
			}, nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		return entity.Secret{}, err
	}

	return entity.Secret{}, fmt.Errorf("unexpected error")
}

func (ps ProcessSecret) ProcessScore(ctx context.Context, score entity.Score) error {
	if err := score.Validate(); err != nil {
		return err
	}

	score.Score = 10 - score.Score

	return ps.repo.InsertUserScore(ctx, score)
}

func (ps ProcessSecret) UpdateUserName(ctx context.Context, score entity.Score) error {
	if err := score.Validate(); err != nil {
		return err
	}
	return ps.repo.UpdateUserName(ctx, score)
}

func (ps ProcessSecret) GetTopPlayer(ctx context.Context) ([]entity.Score, error) {
	return ps.repo.SelectTopPlayer(ctx)
}
