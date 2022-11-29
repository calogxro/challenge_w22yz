package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := r.executeReq(req)

	var m map[string]string
	resp := w.Body.Bytes()
	err := json.Unmarshal(resp, &m)
	message, hasMessage := m["message"]

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err, "json fields should be of type 'string'")
	assert.Equal(t, true, hasMessage, "missing 'message' field in json")
	assert.Equal(t, "pong", message, "message != pong")
}

func TestCreateAnswer(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	w := r.executeReq(req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateNotValidAnswer(t *testing.T) {
	r := NewTestRouter()

	payloads := [][]byte{
		[]byte(`{"key":"", "value": ""}`),
		[]byte(`{"key":"", "value": "John"}`),
		[]byte(`{"key":"name", "value": ""}`),
		[]byte(`{"key":"name"}`),
		[]byte(`{"value": "John"}`),
	}

	for _, jsonStr := range payloads {
		req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(jsonStr))
		w := r.executeReq(req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestCreateExistingAnswer(t *testing.T) {
	r := NewTestRouter()

	var w *httptest.ResponseRecorder

	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
		w = r.executeReq(req)
	}

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestFindAnswer(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	r.executeReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("GET", "/answers/"+key, nil)
	w := r.executeReq(req)

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
	r := NewTestRouter()

	key := testAnswer.Key
	req, _ := http.NewRequest("GET", "/answers/"+key, nil)
	w := r.executeReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAnswer(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	r.executeReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("PATCH", "/answers/"+key, bytes.NewBuffer(testJson))
	w := r.executeReq(req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateNotValidAnswer(t *testing.T) {
	r := NewTestRouter()

	jsonStr := []byte(`{"value": ""}`)
	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(jsonStr))
	w := r.executeReq(req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateNotExistingAnswer(t *testing.T) {
	r := NewTestRouter()

	key := testAnswer.Key
	req, _ := http.NewRequest("PATCH", "/answers/"+key, bytes.NewBuffer(testJson))
	w := r.executeReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteAnswer(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	r.executeReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("DELETE", "/answers/"+key, nil)
	w := r.executeReq(req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteNotExistingAnswer(t *testing.T) {
	r := NewTestRouter()

	key := testAnswer.Key
	req, _ := http.NewRequest("DELETE", "/answers/"+key, nil)
	w := r.executeReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetHistory(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	r.executeReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w := r.executeReq(req)

	var events []Event
	resp := w.Body.Bytes()
	err1 := json.Unmarshal(resp, &events)

	var answer Answer
	event := events[len(events)-1]
	err2 := json.Unmarshal([]byte(event.Data), &answer)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err1)
	assert.Equal(t, 1, len(events))

	assert.True(t, json.Valid(event.Data), "not valid json")
	assert.Nil(t, err2)
	assert.Equal(t, ANSWER_CREATED_EVENT, event.Type)
	assert.Equal(t, testAnswer.Key, answer.Key)
	assert.Equal(t, testAnswer.Value, answer.Value)
}

func TestGetHistoryNotExistingAnswer(t *testing.T) {
	r := NewTestRouter()

	key := testAnswer.Key
	req, _ := http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w := r.executeReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetHistoryDeletedAnswer(t *testing.T) {
	r := NewTestRouter()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	w := r.executeReq(req)

	assert.Equal(t, http.StatusCreated, w.Code)

	key := testAnswer.Key
	req, _ = http.NewRequest("DELETE", "/answers/"+key, nil)
	w = r.executeReq(req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w = r.executeReq(req)

	var events []Event
	resp := w.Body.Bytes()
	err1 := json.Unmarshal(resp, &events)

	var answer Answer
	event := events[len(events)-1]
	err2 := json.Unmarshal([]byte(event.Data), &answer)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err1)
	assert.Equal(t, 2, len(events))

	assert.True(t, json.Valid(event.Data), "not valid json")
	assert.Nil(t, err2)
	assert.Equal(t, ANSWER_DELETED_EVENT, event.Type)
	assert.Equal(t, testAnswer.Key, answer.Key)
	assert.Equal(t, "", answer.Value)
}
