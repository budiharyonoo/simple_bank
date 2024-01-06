package db

import (
	"context"
	"github.com/budiharyonoo/simple_bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: "secret",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	// Check if there is any error
	require.NoError(t, err)

	// Check if user is not empty
	require.NotEmpty(t, user)

	// Check payload from arg and created user are the same
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)

	// Check if password_changed_at and CreatedAt successfuly auto generated
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := createRandomUser(t)
	fetchUser, err := testQueries.GetUser(context.Background(), newUser.Username)

	// Check if there is any error
	require.NoError(t, err)

	// Check if account is not empty
	require.NotEmpty(t, fetchUser)

	// Check if value identical
	require.Equal(t, newUser.HashedPassword, fetchUser.HashedPassword)
	require.Equal(t, newUser.Email, fetchUser.Email)
	require.Equal(t, newUser.FullName, fetchUser.FullName)
	require.WithinDuration(t, newUser.PasswordChangedAt, fetchUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, newUser.CreatedAt, fetchUser.CreatedAt, time.Second)
}
