package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockdb "github.com/tgfukuda/be-master/db/mock"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/token"
	"github.com/tgfukuda/be-master/util"
)

func TestCreateTransfer(t *testing.T) {
	user1, _ := randomUser(t)
	account1 := randomAccount(user1.Username)
	user2, _ := randomUser(t)
	account2 := randomAccount(user2.Username)
	account2.Currency = account1.Currency
	user3, _ := randomUser(t)
	account3 := randomAccount(user3.Username)
	amount := account1.Balance / 2
	transfer := randomTransfer(account1, account2, amount)
	entry1 := randomEntry(account1, -amount)
	entry2 := randomEntry(account2, amount)
	result := db.TransferTxResult{
		Transfer:    transfer,
		FromAccount: account1,
		ToAccount:   account2,
		FromEntry:   entry1,
		ToEntry:     entry2,
	}

	RunTestCases(t, []APITestCase{
		{
			name:   "OK",
			path:   "/transfers",
			method: http.MethodPost,
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{FromAccountID: account1.ID, ToAccountID: account2.ID, Amount: amount})).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				requireMatchTransferTxResult(t, recoder.Body, result)
			},
		},
		{
			name:   "UnAuthorizedFrom",
			path:   "/transfers",
			method: http.MethodPost,
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account3.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name:   "InvalidFromId",
			path:   "/transfers",
			method: http.MethodPost,
			body: gin.H{
				"to_account_id": account2.ID,
				"amount":        amount,
				"currency":      account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InvalidToId",
			path:   "/transfers",
			method: http.MethodPost,
			body: gin.H{
				"from_account_id": account1.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InvalidAmount",
			path:   "/transfers",
			method: http.MethodPost,
			body: gin.H{
				"from_account_id": "",
				"to_account_id":   account2.ID,
				"amount":          0,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:   "InvalidCurrency",
			path:   "/transfers",
			method: http.MethodPost,
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "NOT A CURRENCY",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	})
}

func randomEntry(account db.Account, amount int64) db.Entry {
	return db.Entry{
		ID:        util.RandomInt(1, 1000),
		AccountID: account.ID,
		Amount:    amount,
	}
}

func randomTransfer(from, to db.Account, amount int64) db.Transfer {
	return db.Transfer{
		ID:            util.RandomInt(1, 1000),
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        amount,
	}
}

func requireMatchTransferTxResult(t *testing.T, body *bytes.Buffer, result db.TransferTxResult) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotResult db.TransferTxResult
	err = json.Unmarshal(data, &gotResult)
	assert.NoError(t, err)

	assert.Equal(t, gotResult, result)
}
