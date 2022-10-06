package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var testStore EventStore
var testRouter *gin.Engine

var testAnswer = &Answer{Key: "name", Value: "John"}
var testJson = []byte(`{"key": "name", "value": "John"}`)

func executeTestReq(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func testsSetup() {
	testStore = &MongoStore{
		dbUser: dbUser,
		dbPass: dbPass,
		dbHost: dbHost,
		dbPort: dbPort,
		dbName: dbName,
		dbColl: "test_" + dbColl,
	}
	testCtrl := NewController(testStore)
	testRouter = NewRouter(testCtrl)
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
