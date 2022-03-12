package usecase

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/frisk038/hangman-server/business"
	"github.com/frisk038/hangman-server/business/entity"
)

const maxNbWord = 386264

type repository interface {
	GetYesterdayNumber(ctx context.Context) (int, error)
	InsertTodaySecret(ctx context.Context, secret entity.Secret) error
	SelectTodaySecret(ctx context.Context) (entity.Secret, error)
	InsertUserScore(ctx context.Context, score entity.Score) error
	UpdateUserName(ctx context.Context, score entity.Score) error
	SelectTopPlayer(ctx context.Context, secretNum int) ([]entity.Score, error)
}

type ProcessSecret struct {
	repo repository
}

func NewProcessSecret(repo repository) ProcessSecret {
	return ProcessSecret{repo: repo}
}

func (ps ProcessSecret) InsertSecretTask() {
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

func (ps ProcessSecret) ProcessScore(ctx context.Context, score entity.Score) error {
	if err := score.Validate(); err != nil {
		return err
	}
	return ps.repo.InsertUserScore(ctx, score)
}

func (ps ProcessSecret) UpdateUserName(ctx context.Context, score entity.Score) error {
	if err := score.Validate(); err != nil {
		return err
	}
	score.UserName = strings.ToUpper(score.UserName)
	return ps.repo.UpdateUserName(ctx, score)
}

func (ps ProcessSecret) GetTopPlayer(ctx context.Context, secretNum int) ([]entity.Score, error) {
	if secretNum <= 0 {
		return nil, business.SecretNumNotValid
	}

	return ps.repo.SelectTopPlayer(ctx, secretNum)
}
