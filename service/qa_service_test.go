package service

import (
	"testing"

	es "github.com/calogxro/qaservice/db/event_store"
	rr "github.com/calogxro/qaservice/db/read_repository"

	"github.com/calogxro/qaservice/domain"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	// Setup

	es := es.NewEventStoreStub()
	service := NewQAService(es)
	rr := rr.NewReadRepositoryStub()
	projector := NewProjector(rr)
	projection := NewQAProjection(rr)

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
	assert.IsType(t, &domain.KeyNotFound{}, err)

	// History

	events, _ := service.GetHistory(answer.Key)
	assert.Equal(t, 3, len(events))
}

// sequence allowed:
// create → delete → create → update
func TestCreateDeleteCreateUpdate(t *testing.T) {
	service := NewQAService(es.NewEventStoreStub())

	var testAnswer = domain.Answer{Key: "name", Value: "John"}
	var err error

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.Nil(t, err)
}

// sequence allowed:
// create → update → delete → create → update
func TestCreateUpdateDeleteCreateUpdate(t *testing.T) {
	service := NewQAService(es.NewEventStoreStub())

	var testAnswer = domain.Answer{Key: "name", Value: "John"}
	var err error

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.Nil(t, err)
}

// sequence not allowed:
// create → delete → update
func TestCreateDeleteUpdate(t *testing.T) {
	service := NewQAService(es.NewEventStoreStub())

	var testAnswer = domain.Answer{Key: "name", Value: "John"}
	var err error

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.NotNil(t, err)
}

// sequence not allowed:
// create → create
func TestCreateCreate(t *testing.T) {
	service := NewQAService(es.NewEventStoreStub())

	var testAnswer = domain.Answer{Key: "name", Value: "John"}
	var err error

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.NotNil(t, err)
}

// sequence not allowed:
// create → delete → delete
func TestCreateDeleteDelete(t *testing.T) {
	service := NewQAService(es.NewEventStoreStub())

	var testAnswer = domain.Answer{Key: "name", Value: "John"}
	var err error

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.NotNil(t, err)
}
