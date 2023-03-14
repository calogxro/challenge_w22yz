package http

import (
	"net/http"
	"testing"

	"github.com/calogxro/qaservice/domain"
	rrmemdb "github.com/calogxro/qaservice/projection/repository/memory"
	"github.com/calogxro/qaservice/projection/service/projection"
	"github.com/calogxro/qaservice/shared/testutil"
	"github.com/gin-gonic/gin"
)

var testAnswer = domain.Answer{Key: "name", Value: "John"}

func newTestServer(rr *rrmemdb.ReadRepository) http.Handler {
	//rr := rrmemdb.New()
	service := projection.New(rr)
	handler := MakeHandler(service, gin.New())
	return handler
}

func TestPing(t *testing.T) {
	r := newTestServer(rrmemdb.New())

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

func TestFindAnswer(t *testing.T) {
	rr := rrmemdb.New()
	rr.CreateAnswer(testAnswer)

	r := newTestServer(rr)

	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description:  "test find answer",
		Method:       http.MethodGet,
		Url:          "/answers/" + testAnswer.Key,
		Request:      nil,
		Response:     testAnswer,
		ResponseCode: http.StatusOK,
	}})
}

func TestFindNotExistingAnswer(t *testing.T) {
	rr := rrmemdb.New()

	r := newTestServer(rr)

	testutil.RunServiceTests(t, r, []testutil.TestCase{{
		Description: "test find not-existing answer",
		Method:      http.MethodGet,
		Url:         "/answers/" + testAnswer.Key,
		Request:     nil,
		Response: struct {
			Error string `json:"error"`
		}{
			Error: domain.ErrKeyNotFound.Error(),
		},
		ResponseCode: http.StatusNotFound,
	}})
}
