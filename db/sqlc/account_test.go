package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tgfukuda/be-master/util"
)

func createRandAccount(t *testing.T) Account {
	user := createRandUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Millisecond)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Millisecond)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts2, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, accounts2)

	for _, account := range accounts2 {
		assert.NotEmpty(t, account)
	}
}
