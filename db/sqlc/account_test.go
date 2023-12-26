package db

import (
	"context"
	"database/sql"
	"github.com/budiharyonoo/simple_bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	// Check if there is any error
	require.NoError(t, err)

	// Check if account is not empty
	require.NotEmpty(t, account)

	// Check payload from arg and created account are the same
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// Check if ID and CreatedAt successfuly auto generated
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	fetchAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	// Check if there is any error
	require.NoError(t, err)

	// Check if account is not empty
	require.NotEmpty(t, fetchAccount)

	// Check if value identical
	require.Equal(t, newAccount.ID, fetchAccount.ID)
	require.Equal(t, newAccount.Owner, fetchAccount.Owner)
	require.Equal(t, newAccount.Balance, fetchAccount.Balance)
	require.Equal(t, newAccount.Currency, fetchAccount.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, fetchAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: utils.RandomMoney(),
	}
	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	// Check if there is any error
	require.NoError(t, err)

	// Check if account is not empty
	require.NotEmpty(t, updatedAccount)

	// Check if value identical
	require.Equal(t, arg.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), newAccount.ID)

	// Check if there is any error
	require.NoError(t, err)

	fetchDeletedAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	// Must be error because the data already deleted (Missing from DB)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, fetchDeletedAccount)
}

func TestListAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := testQueries.ListAllAccounts(context.Background())

	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.NotZero(t, account.ID)
	}
}

func TestListAccountsWithPagination(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsWithPaginationParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccountsWithPagination(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.NotZero(t, account.ID)
	}
}
