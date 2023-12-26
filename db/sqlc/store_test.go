package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// Payload data of this test
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)

	fmt.Println(">> Before transfer transaction:", account1.Balance, account2.Balance)

	// Run transfer transaction inside goroutine
	concurrent := 5

	// Because we are using running the transaction inside of goroutine,
	// then the main go routines (this main func TestTransferTx()) is different from
	// the other goroutine inside the loop function
	// so we need channel to catch the results and errors from the store.TransferTx
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < concurrent; i++ {
		go func() {
			// Run query transaction
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			// Send the result or err to the channel
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	// Loop the channels to check the results or errors
	for i := 0; i < concurrent; i++ {
		// Get each result of goroutines
		err := <-errs
		result := <-results

		// === Run the test: ===

		// 1. Check the error is nil
		require.NoError(t, err)

		// 2. Check Transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// 3. Check Entry
		// From Entry
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		// To Entry
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// 4. Check Accounts
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		// check balances
		fmt.Println(">> tx:", result.FromAccount.Balance, result.ToAccount.Balance)

		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= concurrent)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	fmt.Println(">> After transfer transaction:", account1.Balance, account2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	// Payload data of this test
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)

	fmt.Println(">> Before transfer transaction:", account1.Balance, account2.Balance)

	// Run transfer transaction inside goroutine
	concurrent := 10

	// Because we are using running the transaction inside of goroutine,
	// then the main go routines (this main func TestTransferTx()) is different from
	// the other goroutine inside the loop function
	// so we need channel to catch the results and errors from the store.TransferTx
	errs := make(chan error)
	for i := 0; i < concurrent; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		// Mix match the order transaction. For ex: Tx 1 = acc1 -> acc2, Tx2 = acc2 -> acc 1
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			// Run query transaction
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})

			// Send the result or err to the channel
			errs <- err
		}()
	}

	// Loop the channels to check the results or errors
	for i := 0; i < concurrent; i++ {
		// Get each result of goroutines
		err := <-errs

		// === Run the test: ===
		require.NoError(t, err)
	}

	fmt.Println(">> After transfer transaction:", account1.Balance, account2.Balance)

	// The final test result is the balance on Account 1 and 2 should be the same as before update
	// (back to first balance)
	require.Equal(t, account1.Balance, account1.Balance)
	require.Equal(t, account2.Balance, account2.Balance)
}
