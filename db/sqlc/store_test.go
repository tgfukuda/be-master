package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandAccount(t)
	account2 := createRandAccount(t)

	fmt.Printf(">>> before: %d, %d\n", account1.Balance, account2.Balance)

	// run a concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	txResults := make(chan TransferTxResult)

	// the number of thread exited
	exited := make(map[int]bool)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			txResults <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-txResults

		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		assert.NotEmpty(t, transfer)
		assert.Equal(t, account1.ID, transfer.FromAccountID)
		assert.Equal(t, account2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotEmpty(t, transfer.ID)
		assert.NotEmpty(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		assert.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, account1.ID, fromEntry.AccountID)
		assert.Equal(t, amount, -fromEntry.Amount)
		assert.NotEmpty(t, fromEntry.ID)
		assert.NotEmpty(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, account2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotEmpty(t, toEntry.ID)
		assert.NotEmpty(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		// check the account's balance
		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, account2.ID, toAccount.ID)

		fmt.Printf(">>> tx: %d, %d\n", fromAccount.Balance, toAccount.Balance)
		// account1 send `amount` to account2 five times, it will be k * `amount` where k is the number of transactions.
		diff1 := account1.Balance - fromAccount.Balance
		// it will be -k * `amount`
		diff2 := account2.Balance - toAccount.Balance
		assert.Equal(t, diff1, -diff2)

		k := int(diff1 / amount)
		assert.True(t, 1 <= k && k <= n)
		assert.NotContains(t, exited, k)
		exited[k] = true
	}

	// check final updated account
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)
	assert.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Printf(">>> after: %d, %d\n", updatedAccount1.Balance, updatedAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandAccount(t)
	account2 := createRandAccount(t)

	fmt.Printf(">>> before: %d, %d\n", account1.Balance, account2.Balance)

	// run a concurrent transfer transaction
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		var fromAccount, toAccount Account

		// switch sender and receiver. account1 on half and account2 on the other.
		if i%2 == 0 {
			fromAccount = account1
			toAccount = account2
		} else {
			fromAccount = account2
			toAccount = account1
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	// check final updated account. should be the same as before
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1.Balance, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)
	assert.Equal(t, account2.Balance, updatedAccount2.Balance)

	fmt.Printf(">>> after: %d, %d\n", updatedAccount1.Balance, updatedAccount2.Balance)
}
