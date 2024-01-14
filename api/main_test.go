package api

import (
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		JWTSecret:           utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
