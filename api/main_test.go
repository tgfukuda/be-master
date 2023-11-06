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
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/mocks"
	"github.com/tgfukuda/be-master/token"
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
	setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
	buildStubs    func(store *mocks.Store)
	checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker)
}

func (tc *APITestCase) Run(t *testing.T) {
	t.Run(tc.name, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mocks.NewStore(t)
		tc.buildStubs(store)

		server := newTestServer(t, store)

		recorder := httptest.NewRecorder()
		var request *http.Request
		var err error
		request, err = http.NewRequest(tc.method, tc.path, requestJsonBody(t, tc.body))
		assert.NoError(t, err)

		tc.setupAuth(t, request, server.tokenMaker)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder, server.tokenMaker)
	})
}

func RunTestCases(t *testing.T, testCases []APITestCase) {
	for _, tc := range testCases {
		tc.Run(t)
	}
}

func requestJsonBody(t *testing.T, req gin.H) *bytes.Reader {
	b, err := json.Marshal(req)
	assert.NoError(t, err)
	return bytes.NewReader(b)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
