package projection

import (
	"testing"

	"github.com/calogxro/qaservice/domain"
	rrmemdb "github.com/calogxro/qaservice/projection/repository/memory"

	"github.com/stretchr/testify/assert"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

func TestFindAnswer(t *testing.T) {
	repo := rrmemdb.New()
	service := New(repo)

	repo.CreateAnswer(testAnswer)

	answer, err := service.GetAnswer(testAnswer.Key)
	assert.Nil(t, err)
	assert.Equal(t, testAnswer, *answer)
}

func TestFindNotExistingAnswer(t *testing.T) {
	repo := rrmemdb.New()
	service := New(repo)

	answer, err := service.GetAnswer(testAnswer.Key)
	assert.Nil(t, answer)
	assert.ErrorIs(t, err, domain.ErrKeyNotFound)
}
