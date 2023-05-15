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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockdb "github.com/tgfukuda/be-master/db/mock"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/util"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:   "NotFound",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   fmt.Sprintf("/accounts/%d", account.ID),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidId",
			path:   fmt.Sprintf("/accounts/%d", 0),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	})
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	account.Balance = 0

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/accounts",
			method: http.MethodPost,
			body:   CreateAccountRequest{Owner: account.Owner, Currency: account.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{Owner: account.Owner, Currency: account.Currency, Balance: 0})).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchAccount(t, recoder.Body, account)
			},
		},
		{
			name:   "InvalidRequest",
			path:   "/accounts",
			method: http.MethodPost,
			body:   CreateAccountRequest{Owner: "", Currency: account.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{Owner: account.Owner, Currency: account.Currency, Balance: 0})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   "/accounts",
			method: http.MethodPost,
			body:   CreateAccountRequest{Owner: account.Owner, Currency: account.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{Owner: account.Owner, Currency: account.Currency, Balance: 0})).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	})
}

func TestListAccountsAPI(t *testing.T) {
	accounts := []db.Account{randomAccount(), randomAccount(), randomAccount()}

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 10),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 10, Offset: 10})).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchAccounts(t, recoder.Body, accounts)
			},
		},
		{
			name:   "InvalidRequest",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 100),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 10, Offset: 10})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 2, 10),
			method: http.MethodGet,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 10, Offset: 10})).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	})
}

func TestDeleteAccount(t *testing.T) {
	account := randomAccount()

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   DeleteAccountRequest{ID: account.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SafeDeleteAccountTx(gomock.Any(), gomock.Eq(db.SafeDeleteAccountTxParams{ID: account.ID})).
					Times(1).
					Return(db.SafeDeleteAccountTxResult{Account: account}, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchSafeDeleteTxResult(t, recoder.Body, db.SafeDeleteAccountTxResult{Account: account})
			},
		},
		{
			name:   "NotFound",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   DeleteAccountRequest{ID: account.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SafeDeleteAccountTx(gomock.Any(), gomock.Eq(db.SafeDeleteAccountTxParams{ID: account.ID})).
					Times(1).
					Return(db.SafeDeleteAccountTxResult{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:   "InvalidRequest",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   struct{ id string }{}, // missing id field
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SafeDeleteAccountTx(gomock.Any(), gomock.Eq(db.SafeDeleteAccountTxParams{ID: account.ID})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InternalError",
			path:   "/delete_account",
			method: http.MethodPost,
			body:   DeleteAccountRequest{ID: account.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SafeDeleteAccountTx(gomock.Any(), gomock.Eq(db.SafeDeleteAccountTxParams{ID: account.ID})).
					Times(1).
					Return(db.SafeDeleteAccountTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	})
}

func randomAccount() db.Account {
	return db.Account{
		ID:        util.RandomInt(1, 100),
		Owner:     util.RandomOwner(),
		Balance:   util.RandomBalance(),
		Currency:  util.RandomCurrency(),
		CreatedAt: util.RandomDate(),
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

func requireMatchSafeDeleteTxResult(t *testing.T, body *bytes.Buffer, res db.SafeDeleteAccountTxResult) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotRes db.SafeDeleteAccountTxResult
	err = json.Unmarshal(data, &gotRes)
	assert.NoError(t, err)

	assert.Equal(t, gotRes, res)
}
