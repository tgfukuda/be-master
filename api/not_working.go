package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockdb "github.com/tgfukuda/be-master/db/mock"
)

type APITestCase struct {
	name          string
	path          string
	method        string
	body          interface{} // expected json struct. nil if no body.
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
}

func (tc *APITestCase) Run(t *testing.T) {
	t.Run(tc.name, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)

		server := NewServer(store)
		recorder := httptest.NewRecorder()

		var request *http.Request
		var err error
		if tc.body == nil {
			request, err = http.NewRequest(tc.method, tc.path, nil)
		} else {
			request, err = http.NewRequest(tc.method, tc.path, requestJsonBody2(t, tc.body))
		}
		assert.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	})
}

func RunTestCases(t *testing.T, testCases []APITestCase) {
	for _, tc := range testCases {
		tc.Run(t)
	}
}

func requestJsonBody2(t *testing.T, req any) *bytes.Reader {
	b, err := json.Marshal(req)
	assert.NoError(t, err)
	return bytes.NewReader(b)
}
