package entity

import (
	"fmt"
	"strings"

	"github.com/frisk038/hangman-server/business"
	"github.com/gofrs/uuid"
)

// Secret struct holds secret info
type Secret struct {
	SecretWord string
	Number     int
}

// Score holds user score info
type Score struct {
	UserID    uuid.UUID
	SecretNum int
	Score     int
	UserName  string
}

func (s Score) Validate() error {
	if s.SecretNum <= 0 {
		return fmt.Errorf("%w : (%s|%d)", business.SecretNumNotValid, s.UserID.String(), s.SecretNum)
	}
	if s.Score < 0 || s.Score > 10 {
		return fmt.Errorf("%w : (%s|%d)", business.ScoreNotValid, s.UserID.String(), s.Score)
	}
	if len(s.UserName) > 3 {
		return fmt.Errorf("%w: (%s|%s)", business.UsernameNotValid, s.UserID.String(), s.UserName)
	}

	s.UserName = strings.ToUpper(s.UserName)

	return nil
}
