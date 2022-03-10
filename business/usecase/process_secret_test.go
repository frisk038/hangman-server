package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDailySecret(t *testing.T) {
	os.Chdir("../../")
	t.Run("OK", func(t *testing.T) {
		word, err := GenerateDailySecret()
		fmt.Println(word)
		assert.NoError(t, err)
		assert.NotEmpty(t, word)
	})
}
