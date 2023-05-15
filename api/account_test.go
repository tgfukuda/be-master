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

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
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
			name:      "NotFound",
			accountID: account.ID,
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
			name:      "InternalError",
			accountID: account.ID,
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
			name:      "InvalidId",
			accountID: 0,
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
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		params        CreateAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: CreateAccountRequest{Owner: account.Owner, Currency: account.Currency},
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
			params: CreateAccountRequest{Owner: "", Currency: account.Currency},
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
			params: CreateAccountRequest{Owner: account.Owner, Currency: account.Currency},
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, requestJsonBody(t, tc.params))
			assert.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	accounts := []db.Account{randomAccount(), randomAccount(), randomAccount()}

	testCases := []struct {
		name          string
		param         ListAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: ListAccountRequest{PageID: 2, PageSize: 10},
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
			name:  "InvalidRequest",
			param: ListAccountRequest{PageID: 2, PageSize: 100},
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
			name:  "InternalError",
			param: ListAccountRequest{PageID: 2, PageSize: 10},
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.param.PageID, tc.param.PageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		param         any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: DeleteAccountRequest{ID: account.ID},
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
			name:  "NotFound",
			param: DeleteAccountRequest{ID: account.ID},
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
			name:  "InvalidRequest",
			param: struct{ id string }{}, // missing id field
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
			name:  "InternalError",
			param: DeleteAccountRequest{ID: account.ID},
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/delete_account"
			body := requestJsonBody(t, tc.param)
			request, err := http.NewRequest(http.MethodPost, url, body)
			assert.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
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

func requestJsonBody(t *testing.T, req any) *bytes.Reader {
	b, err := json.Marshal(req)
	assert.NoError(t, err)
	return bytes.NewReader(b)
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
