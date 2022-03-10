package usecase

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/frisk038/hangman-server/business/entity"
)

const maxNbWord = 386264

func GenerateDailySecret() (entity.Secret, error) {
	f, err := os.Open("./dico.txt")
	if err != nil {
		return entity.Secret{}, err
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
