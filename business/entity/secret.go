package entity

import "github.com/gofrs/uuid"

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
}
