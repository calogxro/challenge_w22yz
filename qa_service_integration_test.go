package main

import (
	"testing"

	Ω "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestServiceWithMySQL(t *testing.T) {
	db, _ := initMySQL()
	db.Exec("DELETE FROM answer")

	// Setup

	es := NewEventStoreStub()
	service := NewQAService(es)
	rr := NewMySQLReadRepository()
	projector := NewProjector(rr)
	projection := NewQAProjection(rr)

	es.Subscribe(func(event *Event) {
		projector.Project(event)
	})

	//Create

	answer := testAnswer
	service.CreateAnswer(answer)
	projAnswer, _ := projection.GetAnswer("name")

	assert.Equal(t, &answer, projAnswer)

	// Update

	answer = Answer{Key: answer.Key, Value: answer.Value + "_2"}
	service.UpdateAnswer(answer)
	projAnswer, _ = projection.GetAnswer("name")

	assert.Equal(t, &answer, projAnswer)

	// Delete

	service.DeleteAnswer(answer.Key)
	projAnswer, err := projection.GetAnswer("name")

	assert.Nil(t, projAnswer)
	assert.NotNil(t, err)
	assert.IsType(t, &KeyNotFound{}, err)

	// History

	events, _ := service.GetHistory(answer.Key)
	assert.Equal(t, 3, len(events))
}

func TestServiceWithEventStoreDB(t *testing.T) {
	g := Ω.NewGomegaWithT(t)

	db, _ := initMySQL()
	db.Exec("DELETE FROM answer")

	// Setup

	es := NewEventStoreDB()
	service := NewQAService(es)
	rr := NewMySQLReadRepository()
	projector := NewProjector(rr)
	projection := NewQAProjection(rr)

	es.deleteStream()

	go es.Subscribe(func(event *Event) {
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
				<*main.Answer | 0xc000382920>: {Key: "", Value: ""}
			to equal
				<*main.Answer | 0xc0003c7fa0>: {Key: "name", Value: "John"}
		FAIL
	*/

	g.Eventually(func() *Answer {
		projAnswer, _ := projection.GetAnswer("name")
		return projAnswer
	}).Should(Ω.Equal(&answer))

	// Update

	answer = Answer{Key: answer.Key, Value: answer.Value + "_2"}
	service.UpdateAnswer(answer)

	g.Eventually(func() *Answer {
		projAnswer, _ := projection.GetAnswer("name")
		return projAnswer
	}).Should(Ω.Equal(&answer))

	// Delete

	service.DeleteAnswer(answer.Key)

	g.Eventually(func() error {
		_, err := projection.GetAnswer("name")
		return err
	}).Should(Ω.Equal(&KeyNotFound{}))
}
