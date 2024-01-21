package token

import (
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker()
	require.NoError(t, err)

	username := utils.RandomOwner()
	tokenDuration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(tokenDuration)

	// Create token
	tokenString, err := maker.CreateToken(username, tokenDuration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Verify the created token
	payload, err := maker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
