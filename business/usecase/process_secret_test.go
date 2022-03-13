package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// type mockRepo struct {
// 	mock.Mock
// }

// func (m mockRepo) getYesterdayNumber(ctx context.Context) (int, error) {
// 	c := m.Called(ctx)
// }
// func (m mockRepo) InsertTodaySecret(ctx context.Context, secret entity.Secret) error
// func (m mockRepo) SelectTodaySecret(ctx context.Context) (entity.Secret, error)

func TestGenerateDailySecret(t *testing.T) {
	os.Chdir("../../")
	t.Run("OK", func(t *testing.T) {

		p := NewProcessSecret(nil)
		word, err := p.generateDailySecret()
		fmt.Println(word)
		assert.NoError(t, err)
		assert.NotEmpty(t, word)
	})
}
