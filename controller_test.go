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
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := executeTestReq(req)

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
	testStore.drop()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	w := executeTestReq(req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateNotValidAnswer(t *testing.T) {
	testStore.drop()

	payloads := [][]byte{
		[]byte(`{"key":"", "value": ""}`),
		[]byte(`{"key":"", "value": "John"}`),
		[]byte(`{"key":"name", "value": ""}`),
		[]byte(`{"key":"name"}`),
		[]byte(`{"value": "John"}`),
	}

	for _, jsonStr := range payloads {
		req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(jsonStr))
		w := executeTestReq(req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestCreateExistingAnswer(t *testing.T) {
	testStore.drop()

	var w *httptest.ResponseRecorder

	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
		w = executeTestReq(req)
	}

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestFindAnswer(t *testing.T) {
	testStore.drop()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	executeTestReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("GET", "/answers/"+key, nil)
	w := executeTestReq(req)

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
	testStore.drop()

	key := testAnswer.Key
	req, _ := http.NewRequest("GET", "/answers/"+key, nil)
	w := executeTestReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAnswer(t *testing.T) {
	testStore.drop()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	executeTestReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("PATCH", "/answers/"+key, bytes.NewBuffer(testJson))
	w := executeTestReq(req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateNotValidAnswer(t *testing.T) {
	testStore.drop()

	jsonStr := []byte(`{"value": ""}`)
	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(jsonStr))
	w := executeTestReq(req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateNotExistingAnswer(t *testing.T) {
	testStore.drop()

	key := testAnswer.Key
	req, _ := http.NewRequest("PATCH", "/answers/"+key, bytes.NewBuffer(testJson))
	w := executeTestReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteAnswer(t *testing.T) {
	testStore.drop()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	executeTestReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("DELETE", "/answers/"+key, nil)
	w := executeTestReq(req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteNotExistingAnswer(t *testing.T) {
	testStore.drop()

	key := testAnswer.Key
	req, _ := http.NewRequest("DELETE", "/answers/"+key, nil)
	w := executeTestReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetHistory(t *testing.T) {
	testStore.drop()

	req, _ := http.NewRequest("POST", "/answers", bytes.NewBuffer(testJson))
	executeTestReq(req)

	key := testAnswer.Key
	req, _ = http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w := executeTestReq(req)

	var events []interface{}
	resp := w.Body.Bytes()
	err := json.Unmarshal(resp, &events)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, json.Valid(resp), "not valid json")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))

	event, ok1 := events[0].(map[string]interface{})
	answer, ok2 := event["data"].(map[string]interface{})

	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, CREATED_EVENT, event["event"])
	assert.Equal(t, testAnswer.Key, answer["key"])
	assert.Equal(t, testAnswer.Value, answer["value"])
}

func TestGetHistoryNotExistingAnswer(t *testing.T) {
	testStore.drop()

	key := testAnswer.Key
	req, _ := http.NewRequest("GET", "/answers/"+key+"/history", nil)
	w := executeTestReq(req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
