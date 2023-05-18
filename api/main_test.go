package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq" // importing with name _ is special import to tell go not to remove this deps
	"github.com/stretchr/testify/assert"
	mockdb "github.com/tgfukuda/be-master/db/mock"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	assert.NoError(t, err)

	return server
}

type APITestCase struct {
	name          string
	path          string
	method        string
	body          gin.H // expected json struct. nil if no body.
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
}

func (tc *APITestCase) Run(t *testing.T) {
	t.Run(tc.name, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)

		server := newTestServer(t, store)

		recorder := httptest.NewRecorder()
		var request *http.Request
		var err error
		request, err = http.NewRequest(tc.method, tc.path, requestJsonBody(t, tc.body))
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

func requestJsonBody(t *testing.T, req any) *bytes.Reader {
	b, err := json.Marshal(req)
	assert.NoError(t, err)
	return bytes.NewReader(b)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
