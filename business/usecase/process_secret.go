package usecase

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/frisk038/hangman-server/business/entity"
)

const maxNbWord = 386264

func GenerateDailySecret() (entity.Secret, error) {
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
