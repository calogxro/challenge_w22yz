package memory

import (
	"testing"

	"github.com/calogxro/qaservice/domain"
	esmemdb "github.com/calogxro/qaservice/eventstore/repository/memory"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	rrmemdb "github.com/calogxro/qaservice/projection/repository/memory"
	"github.com/calogxro/qaservice/projection/service/projection"
	"github.com/calogxro/qaservice/projector/service/projector"
	"github.com/stretchr/testify/assert"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

func TestService(t *testing.T) {
	// Setup

	es := esmemdb.New()
	rr := rrmemdb.New()

	service := eventstore.New(es)
	projection := projection.New(rr)
	projector := projector.New(rr)

	es.Subscribe(func(event *domain.Event) {
		projector.Project(event)
	})

	//Create

	answer := testAnswer
	service.CreateAnswer(answer)
	projAnswer, _ := projection.GetAnswer("name")

	assert.Equal(t, &answer, projAnswer)

	// Update

	answer = domain.Answer{Key: answer.Key, Value: answer.Value + "_2"}
	service.UpdateAnswer(answer)
	projAnswer, _ = projection.GetAnswer("name")

	assert.Equal(t, &answer, projAnswer)

	// Delete

	service.DeleteAnswer(answer.Key)
	projAnswer, err := projection.GetAnswer("name")

	assert.Nil(t, projAnswer)
	assert.NotNil(t, err)
	assert.IsType(t, domain.ErrKeyNotFound, err)

	// History

	events, _ := service.GetHistory(answer.Key)
	assert.Equal(t, 3, len(events))
}
