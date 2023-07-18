package db

import (
	"testing"

	"github.com/calogxro/qaservice/eventstore/repository/esdb"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	projectionRepo "github.com/calogxro/qaservice/projection/repository/mongodb"
	"github.com/calogxro/qaservice/projection/service/projection"
	gateway "github.com/calogxro/qaservice/projector/gateway/esdbgw"
	projectorRepo "github.com/calogxro/qaservice/projector/repository/mongodb"
	"github.com/calogxro/qaservice/projector/service/projector"

	"github.com/calogxro/qaservice/domain"
	Ω "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

func TestServiceWithDBs(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	g := Ω.NewGomegaWithT(t)

	// Setup

	esRepo := esdb.New()
	projectionRepo := projectionRepo.New()
	projectorRepo := projectorRepo.New()

	eventstore := eventstore.New(esRepo)
	projection := projection.New(projectionRepo)
	projector := projector.New(projectorRepo)

	esRepo.DeleteStream()
	projectorRepo.DeleteAllAnswers()

	esdbgw := gateway.New()
	go esdbgw.Subscribe(func(event *domain.Event) {
		projector.Project(event)
	})

	//Create

	answer := testAnswer
	_, err := eventstore.CreateAnswer(answer)

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
	}, 2).Should(Ω.Equal(&answer)) // timeout: 2 sec

	// Update

	answer = domain.Answer{Key: answer.Key, Value: answer.Value + "_2"}
	eventstore.UpdateAnswer(answer)

	g.Eventually(func() *domain.Answer {
		projAnswer, _ := projection.GetAnswer("name")
		return projAnswer
	}).Should(Ω.Equal(&answer))

	// Delete

	eventstore.DeleteAnswer(answer.Key)

	g.Eventually(func() error {
		_, err := projection.GetAnswer("name")
		return err
	}).Should(Ω.Equal(domain.ErrKeyNotFound))
}
