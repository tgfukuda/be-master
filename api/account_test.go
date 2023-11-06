package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/mocks"
	"github.com/tgfukuda/be-master/token"
	"github.com/tgfukuda/be-master/util"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				// build stubs
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:   "UnAuthorizedUser",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				// build stubs
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "NoAuthorization",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				// build stubs
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				// build stubs
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidId",
			path:   fmt.Sprintf("/accounts/%d", 0),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	})
}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	account.Balance = 0

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/accounts",
			method: http.MethodPost,
			body:   gin.H{"currency": account.Currency},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					CreateAccount(mock.Anything, db.CreateAccountParams{Owner: account.Owner, Currency: account.Currency, Balance: 0}).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchAccount(t, recoder.Body, account)
			},
		},
		{
			name:   "NoAuthorization",
			path:   "/accounts",
			method: http.MethodPost,
			body:   gin.H{"currency": "NOT A CURRENCY"},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name:   "InvalidRequest",
			path:   "/accounts",
			method: http.MethodPost,
			body:   gin.H{"currency": "NOT A CURRENCY"},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   "/accounts",
			method: http.MethodPost,
			body:   gin.H{"currency": account.Currency},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					CreateAccount(mock.Anything, db.CreateAccountParams{Owner: account.Owner, Currency: account.Currency, Balance: 0}).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		// TODO: forbidden case
	})
}

func TestListAccountsAPI(t *testing.T) {
	user, _ := randomUser(t)
	n := 5

	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 10),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					ListAccounts(mock.Anything, db.ListAccountsParams{Owner: user.Username, Limit: 10, Offset: 10}).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchAccounts(t, recoder.Body, accounts)
			},
		},
		{
			name:   "NoAuthorization",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 10),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name:   "InvalidRequest",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 100),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 10),
			method: http.MethodGet,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					ListAccounts(mock.Anything, db.ListAccountsParams{Owner: user.Username, Limit: 10, Offset: 10}).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	})
}

func TestDeleteAccount(t *testing.T) {
	user, _ := randomUser(t)
	attacker, _ := randomUser(t)
	account := randomAccount(user.Username)

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   gin.H{"id": account.ID},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(account, nil)

				store.EXPECT().
					DeleteAccount(mock.Anything, account.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchSafeDeleteTxResult(t, recoder.Body, account)
			},
		},
		{
			name:   "InvalidRequest",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "NotFound",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   gin.H{"id": account.ID},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:   "InternalErrorSQLGetReq",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   gin.H{"id": account.ID},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:   "InternalErrorSQLDeleteReq",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   gin.H{"id": account.ID},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(account, nil)

				store.EXPECT().
					DeleteAccount(mock.Anything, account.ID).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:   "Unauthorized",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   gin.H{"id": account.ID},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, attacker.Username, time.Minute)
			},
			buildStubs: func(store *mocks.Store) {
				store.EXPECT().
					GetAccount(mock.Anything, account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
	})
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 100),
		Owner:    owner,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requireMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	assert.NoError(t, err)

	assert.Equal(t, gotAccount, account)
}

func requireMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	assert.NoError(t, err)

	assert.Equal(t, gotAccounts, accounts)
}

func requireMatchSafeDeleteTxResult(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotRes map[string]db.Account
	err = json.Unmarshal(data, &gotRes)
	assert.NoError(t, err)

	deleted := gotRes["deleted"]

	assert.Equal(t, deleted, account)
}
