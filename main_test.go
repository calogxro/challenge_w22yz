package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var testStore IEventStore
var testRouter *gin.Engine

var testAnswer = &Answer{Key: "name", Value: "John"}
var testJson = []byte(`{"key": "name", "value": "John"}`)

func executeTestReq(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func testsSetup() {
	testColl := db.Database(DB_NAME).Collection(DB_COLL_TEST)
	testStore = &EventStore{testColl}
	testCtrl := Controller{testStore}
	testRouter = NewRouter(&testCtrl)
}

func TestMain(m *testing.M) {
	// Do stuff BEFORE the tests
	testsSetup()

	// Run the tests
	code := m.Run()

	// Do stuff AFTER the tests
	testStore.drop()

	os.Exit(code)
}
