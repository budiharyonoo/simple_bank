package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// renewAccessTokenRequest used for input payload request and validation
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// renewAccessTokenResponse is for success API response
type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	// Bind JSON payload to the createUserRequest struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	// Validate refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	fmt.Println(refreshPayload)

	// Validate if refresh token session id exists
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		err := errors.New("no session id exists")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := errors.New("your session blocked")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// Validate refresh token
	if session.RefreshToken != req.RefreshToken {
		err := errors.New("refresh token not valid with your current session")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// Validate if refresh token already expired
	if time.Now().After(session.ExpiredAt) {
		err := errors.New("refresh token expired")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// Create new Auth Token
	tokenString, accessTokenPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username, server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          tokenString,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	})
}
