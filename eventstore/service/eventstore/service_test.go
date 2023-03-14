package eventstore

import (
	"encoding/json"
	"testing"

	"github.com/calogxro/qaservice/domain"
	esmemdb "github.com/calogxro/qaservice/eventstore/repository/memory"

	"github.com/stretchr/testify/assert"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

//var testJson = []byte(`{"key": "name", "value": "John"}`)

func TestCreateAnswer(t *testing.T) {
	repo := esmemdb.New()
	service := New(repo)

	event, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	var answer domain.Answer
	err = json.Unmarshal(event.Data, &answer)
	assert.Nil(t, err)
	assert.Equal(t, testAnswer, answer)

	events, err := repo.GetEvents()
	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Equal(t, event, events[len(events)-1])
}

func TestCreateNotValidAnswer(t *testing.T) {
	service := New(esmemdb.New())

	answers := []domain.Answer{
		{Key: "", Value: ""},
		{Key: "", Value: "John"},
		{Key: "name", Value: ""},
		{Key: "name"},
		{Value: "John"},
	}

	for _, testAnswer := range answers {
		_, err := service.CreateAnswer(testAnswer)
		assert.ErrorIs(t, err, domain.ErrInputNotValid)
	}
}

func TestCreateExistingAnswer(t *testing.T) {
	service := New(esmemdb.New())

	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	_, err = service.CreateAnswer(testAnswer)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, domain.ErrKeyExists)
}

func TestUpdateAnswer(t *testing.T) {
	repo := esmemdb.New()
	service := New(repo)

	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	updatedAnswer := domain.Answer{
		Key:   testAnswer.Key,
		Value: testAnswer.Value + "X",
	}
	event, err := service.UpdateAnswer(updatedAnswer)
	assert.Nil(t, err)

	var answer domain.Answer
	err = json.Unmarshal(event.Data, &answer)
	assert.Nil(t, err)
	assert.Equal(t, updatedAnswer, answer)

	events, err := repo.GetEvents()
	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Equal(t, event, events[len(events)-1])
}

func TestUpdateNotValidAnswer(t *testing.T) {
	service := New(esmemdb.New())

	inputs := []domain.Answer{
		{Key: "", Value: ""},
		{Key: "", Value: "John"},
		{Key: "name", Value: ""},
	}

	for _, answer := range inputs {
		_, err := service.UpdateAnswer(answer)
		assert.ErrorIs(t, err, domain.ErrInputNotValid)
	}
}

func TestUpdateNotExistingAnswer(t *testing.T) {
	service := New(esmemdb.New())

	_, err := service.UpdateAnswer(testAnswer)
	assert.ErrorIs(t, err, domain.ErrKeyNotFound)
}

func TestDeleteAnswer(t *testing.T) {
	service := New(esmemdb.New())

	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)
}

func TestDeleteNotExistingAnswer(t *testing.T) {
	service := New(esmemdb.New())

	_, err := service.DeleteAnswer(testAnswer.Key)
	assert.ErrorIs(t, err, domain.ErrKeyNotFound)
}

func TestGetHistory(t *testing.T) {
	service := New(esmemdb.New())

	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	updatedAnswer := domain.Answer{
		Key:   testAnswer.Key,
		Value: testAnswer.Value + "X",
	}
	_, err = service.UpdateAnswer(updatedAnswer)
	assert.Nil(t, err)

	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	events, err := service.GetHistory(testAnswer.Key)
	assert.Nil(t, err)
	assert.NotEmpty(t, events)

	var answer domain.Answer
	event := events[len(events)-3]
	err = json.Unmarshal(event.Data, &answer)
	assert.Nil(t, err)
	assert.Equal(t, testAnswer, answer)
	assert.Equal(t, domain.ANSWER_CREATED_EVENT, event.Type)

	event = events[len(events)-2]
	assert.Equal(t, domain.ANSWER_UPDATED_EVENT, event.Type)

	event = events[len(events)-1]
	assert.Equal(t, domain.ANSWER_DELETED_EVENT, event.Type)
}

func TestGetHistoryNotExistingAnswer(t *testing.T) {
	service := New(esmemdb.New())

	_, err := service.GetHistory(testAnswer.Key)
	assert.ErrorIs(t, err, domain.ErrKeyNotFound)
}

/*
func TestGetHistoryDeletedAnswer(t *testing.T) {
	r := router.NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	w := r.ExecuteReq(req)

	assert.Equal(t, http.StatusCreated, w.Code)

	key := testAnswer.Key
	req, _ = http.NewRequest("DELETE", "/answers/"+key, nil)
	w = r.ExecuteReq(req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w = r.ExecuteReq(req)

	var events []domain.Event
	resp := w.Body.Bytes()
	err1 := json.Unmarshal(resp, &events)

	var answer domain.Answer
	event := events[len(events)-1]
	err2 := json.Unmarshal([]byte(event.Data), &answer)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err1)
	assert.Equal(t, 2, len(events))

	assert.True(t, json.Valid(event.Data), "not valid json")
	assert.Nil(t, err2)
	assert.Equal(t, domain.ANSWER_DELETED_EVENT, event.Type)
	assert.Equal(t, testAnswer.Key, answer.Key)
	assert.Equal(t, "", answer.Value)
}
*/
/*

func TestFindAnswer(t *testing.T) {
	r := router.NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	r.ExecuteReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("GET", "/answers/"+key, nil)
	w := r.ExecuteReq(req)

	var m map[string]string
	resp := w.Body.Bytes()
	err := json.Unmarshal(resp, &m)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(m))
	assert.Equal(t, testAnswer.Key, m["key"])
	assert.Equal(t, testAnswer.Value, m["value"])
}

func TestFindNotExistingAnswer(t *testing.T) {
	r := router.NewTestRouter()

	key := testAnswer.Key
	req, _ := http.NewRequest("GET", "/answers/"+key, nil)
	w := r.ExecuteReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}




func TestGetHistoryDeletedAnswer(t *testing.T) {
	r := router.NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	w := r.ExecuteReq(req)

	assert.Equal(t, http.StatusCreated, w.Code)

	key := testAnswer.Key
	req, _ = http.NewRequest("DELETE", "/answers/"+key, nil)
	w = r.ExecuteReq(req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w = r.ExecuteReq(req)

	var events []domain.Event
	resp := w.Body.Bytes()
	err1 := json.Unmarshal(resp, &events)

	var answer domain.Answer
	event := events[len(events)-1]
	err2 := json.Unmarshal([]byte(event.Data), &answer)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err1)
	assert.Equal(t, 2, len(events))

	assert.True(t, json.Valid(event.Data), "not valid json")
	assert.Nil(t, err2)
	assert.Equal(t, domain.ANSWER_DELETED_EVENT, event.Type)
	assert.Equal(t, testAnswer.Key, answer.Key)
	assert.Equal(t, "", answer.Value)
}
*/
