package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDailySecret(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		word, err := GenerateDailySecret()
		assert.NoError(t, err)
		assert.NotEmpty(t, word)
	})
}
