package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

var testAnswer = Answer{Key: "name", Value: "John"}
var testJson = []byte(`{"key": "name", "value": "John"}`)

type TestRouter struct {
	router *gin.Engine
}

func NewTestRouter() *TestRouter {
	es := NewEventStoreStub()
	rr := NewReadRepositoryStub()
	service := NewQAService(es)
	projection := NewQAProjection(rr)
	ctrl := NewController(service, projection)
	router := NewRouter(ctrl)

	projector := NewProjector(rr)
	es.Subscribe(func(event *Event) {
		projector.Project(event)
	})

	return &TestRouter{
		router: router,
	}
}

func (r *TestRouter) executeReq(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.router.ServeHTTP(w, req)
	return w
}
