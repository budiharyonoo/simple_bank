package utils

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	// === Positive Case ===
	// Generate Hashed Password
	password := RandomString(8)
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	// Validate the Hashed Password
	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err)

	// Check if you hash the password more than 1 times, the result must be different
	// from the first hashed password
	newHashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, newHashedPassword)
	require.NotEqual(t, hashedPassword, newHashedPassword)

	// === Negative Case ===
	wrongPassword := RandomString(8)
	err = CheckPassword(wrongPassword, password)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
