package http

import (
	"net/http"
	"testing"

	"github.com/calogxro/qaservice/domain"
	esmemdb "github.com/calogxro/qaservice/eventstore/repository/memory"
	"github.com/calogxro/qaservice/eventstore/service/eventstore"
	"github.com/calogxro/qaservice/shared/testutil"
	"github.com/gin-gonic/gin"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

func newTestServer() http.Handler {
	es := esmemdb.New()
	service := eventstore.New(es)
	handler := MakeHandler(service, gin.New())
	return handler
}

func TestPing(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description: "test ping",
		Method:      http.MethodGet,
		Url:         "/ping",
		Request:     nil,
		Response: struct {
			Message string `json:"message"`
		}{
			Message: "pong",
		},
		ResponseCode: http.StatusOK,
	}})
}

func TestCreateAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}})
}

func TestCreateNotValidAnswer(t *testing.T) {
	r := newTestServer()

	testAnswers := []domain.Answer{
		{Key: "", Value: ""},
		{Key: "", Value: "John"},
		{Key: "name", Value: ""},
		{Key: "name"},
		{Value: "John"},
	}

	testCases := []testutil.TestCase{}

	for _, testAnswer := range testAnswers {
		testCases = append(testCases, testutil.TestCase{
			Description:  "test create not-valid answer",
			Method:       http.MethodPost,
			Url:          "/answers",
			Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
			Response:     nil,
			ResponseCode: http.StatusBadRequest,
		})
	}

	testutil.RunServiceTests(t, r, testCases)
	// r := newTestServer()

	// payloads := [][]byte{
	// 	[]byte(`{"key":"", "value": ""}`),
	// 	[]byte(`{"key":"", "value": "John"}`),
	// 	[]byte(`{"key":"name", "value": ""}`),
	// 	[]byte(`{"key":"name"}`),
	// 	[]byte(`{"value": "John"}`),
	// }

	// for _, jsonStr := range payloads {
	// 	body := bytes.NewBuffer(jsonStr)
	// 	req, _ := http.NewRequest("POST", "/answers", body)

	// 	w := httptest.NewRecorder()
	// 	r.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// }
}

func TestCreateExistingAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}, {
		Description:  "test re-create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: domain.ErrKeyExists.Error()},
		ResponseCode: http.StatusConflict,
	}})
}

func TestUpdateAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}, {
		Description:  "test update answer",
		Method:       http.MethodPatch,
		Url:          "/answers/" + testAnswer.Key,
		Request:      UpdateRequest{testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusOK,
	}})
}

func TestUpdateNotValidAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}, {
		Description:  "test update not-valid answer",
		Method:       http.MethodPatch,
		Url:          "/answers/" + testAnswer.Key,
		Request:      CreateRequest{testAnswer.Key, ""},
		Response:     nil,
		ResponseCode: http.StatusBadRequest,
	}})
}

func TestUpdateNotExistingAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test update not-existing answer",
		Method:       http.MethodPatch,
		Url:          "/answers/" + testAnswer.Key,
		Request:      CreateRequest{testAnswer.Key, ""},
		Response:     nil,
		ResponseCode: http.StatusBadRequest,
	}})
}

func TestDeleteAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}, {
		Description:  "test delete answer",
		Method:       http.MethodDelete,
		Url:          "/answers/" + testAnswer.Key,
		Request:      nil,
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusOK,
	}})
}

func TestDeleteNotExistingAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test delete not-existing answer",
		Method:       http.MethodDelete,
		Url:          "/answers/" + testAnswer.Key,
		Request:      nil,
		Response:     Response{Key: testAnswer.Key, Err: domain.ErrKeyNotFound.Error()},
		ResponseCode: http.StatusNotFound,
	}})
}

func TestGetHistoryOfNotExistingAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test get history of not-existing answer",
		Method:       http.MethodGet,
		Url:          "/answers/" + testAnswer.Key + "/history",
		Request:      nil,
		Response:     Response{Key: testAnswer.Key, Err: domain.ErrKeyNotFound.Error()},
		ResponseCode: http.StatusNotFound,
	}})
}

func TestGetHistory(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}, {
		Description: "test get history of answer",
		Method:      http.MethodGet,
		Url:         "/answers/" + testAnswer.Key + "/history",
		Request:     nil,
		Response: Response{Key: testAnswer.Key, Err: "", Data: []event{{
			Type: domain.ANSWER_CREATED_EVENT,
			Data: testAnswer,
		}}},
		ResponseCode: http.StatusOK,
	}})
}

func TestGetHistoryDeletedAnswer(t *testing.T) {
	r := newTestServer()
	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test create answer",
		Method:       http.MethodPost,
		Url:          "/answers",
		Request:      CreateRequest{testAnswer.Key, testAnswer.Value},
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusCreated,
	}, {
		Description:  "test delete answer",
		Method:       http.MethodDelete,
		Url:          "/answers/" + testAnswer.Key,
		Request:      nil,
		Response:     Response{Key: testAnswer.Key, Err: ""},
		ResponseCode: http.StatusOK,
	}, {
		Description: "test get history of deleted answer",
		Method:      http.MethodGet,
		Url:         "/answers/" + testAnswer.Key + "/history",
		Request:     nil,
		Response: Response{Key: testAnswer.Key, Err: "", Data: []event{{
			Type: domain.ANSWER_CREATED_EVENT,
			Data: testAnswer,
		}, {
			Type: domain.ANSWER_DELETED_EVENT,
			Data: domain.Answer{Key: testAnswer.Key},
		}}},
		ResponseCode: http.StatusOK,
	}})
}
