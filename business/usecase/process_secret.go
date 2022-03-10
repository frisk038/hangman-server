package usecase

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/frisk038/hangman-server/business/entity"
)

const maxNbWord = 386264

type repository interface {
	GetYesterdayNumber(ctx context.Context) (int, error)
	InsertTodaySecret(ctx context.Context, secret entity.Secret) error
	SelectTodaySecret(ctx context.Context) (entity.Secret, error)
}

type ProcessSecret struct {
	repo repository
}

func NewProcessSecret(repo repository) ProcessSecret {
	return ProcessSecret{repo: repo}
}

func (ps ProcessSecret) InsertSecret(ctx context.Context) error {
	secret, err := ps.generateDailySecret()
	if err != nil {
		return err
	}

	nb, err := ps.repo.GetYesterdayNumber(ctx)
	if err != nil {
		return err
	}

	secret.Number = nb + 1
	return ps.repo.InsertTodaySecret(ctx, secret)
}

func (ps ProcessSecret) GetSecret(ctx context.Context) (entity.Secret, error) {
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
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	wantedLine := random.Intn(maxNbWord)
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
