package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandAccount(t)
	account2 := createRandAccount(t)

	// run a concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	txResults := make(chan TransferTxResult)

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
	}
}
