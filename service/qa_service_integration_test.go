package service

import (
	"testing"

	es "github.com/calogxro/qaservice/db/event_store"
	rr "github.com/calogxro/qaservice/db/read_repository"

	"github.com/calogxro/qaservice/domain"
	Ω "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

func TestServiceWithDBs(t *testing.T) {
	g := Ω.NewGomegaWithT(t)

	// Setup

	es := es.NewEventStoreDB()
	service := NewQAService(es)
	rr := rr.NewMongoRepository()
	projector := NewProjector(rr)
	projection := NewQAProjection(rr)

	es.DeleteStream()
	rr.DeleteAllAnswers()

	go es.Subscribe(func(event *domain.Event) {
		projector.Project(event)
	})

	//Create

	answer := testAnswer
	_, err := service.CreateAnswer(answer)

	assert.Nil(t, err)

	// projAnswer, _ := projection.GetAnswer("name")
	// assert.Equal(t, &answer, projAnswer)
	/*
		--- FAIL: TestServiceWithEventStoreDB (1.22s)
		qa_service_2_test.go:123:
			Timed out after 1.001s.
			Expected
				<*main.domain.Answer | 0xc000382920>: {Key: "", Value: ""}
			to equal
				<*main.domain.Answer | 0xc0003c7fa0>: {Key: "name", Value: "John"}
		FAIL
	*/

	g.Eventually(func() *domain.Answer {
		projAnswer, _ := projection.GetAnswer("name")
		return projAnswer
	}).Should(Ω.Equal(&answer))

	// Update

	answer = domain.Answer{Key: answer.Key, Value: answer.Value + "_2"}
	service.UpdateAnswer(answer)

	g.Eventually(func() *domain.Answer {
		projAnswer, _ := projection.GetAnswer("name")
		return projAnswer
	}).Should(Ω.Equal(&answer))

	// Delete

	service.DeleteAnswer(answer.Key)

	g.Eventually(func() error {
		_, err := projection.GetAnswer("name")
		return err
	}).Should(Ω.Equal(&domain.KeyNotFound{}))
}
