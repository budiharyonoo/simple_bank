package api

import (
	"errors"
	"github.com/budiharyonoo/simple_bank/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authHeaderKey  = "authorization"
	authType       = "bearer"
	authPayloadKey = "auth_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Missing authorization key in header
		authHeader := ctx.GetHeader(authHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("no authorization header exists")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Split authHeader string by space
		splitAuth := strings.Fields(authHeader)

		// Missing auth type or token
		if len(splitAuth) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Invalid auth token type, should be "Bearer"
		if authType != strings.ToLower(splitAuth[0]) {
			err := errors.New("invalid authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(splitAuth[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Append token payload to the request
		ctx.Set(authPayloadKey, payload)

		ctx.Next()
	}
}
