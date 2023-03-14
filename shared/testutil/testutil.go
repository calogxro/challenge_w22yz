package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Description  string
	Method       string
	Url          string
	Request      interface{} // can be nil
	Response     interface{} // nothing asserted if nil
	ResponseCode int
}

func RunServiceTests(t *testing.T, r http.Handler, tests []TestCase) {
	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			runServiceTest(t, r, test)
		})
	}
}

func runServiceTest(t *testing.T, r http.Handler, test TestCase) {
	var body io.Reader
	if test.Request != nil {
		reqBytes, err := json.Marshal(test.Request)
		assert.Nil(t, err)
		body = bytes.NewBuffer(reqBytes)
	}

	req, _ := http.NewRequest(test.Method, test.Url, body)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Body.Bytes()

	// nothing asserted if nil
	if test.Response != nil {
		assert.True(t, json.Valid(resp), "not valid json")
		respWanted, err := json.Marshal(test.Response)
		assert.Nil(t, err)
		assert.Equal(t, respWanted, resp)
	}

	assert.Equal(t, test.ResponseCode, w.Code)
}
