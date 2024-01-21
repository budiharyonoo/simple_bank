package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload is data that store in JWT & PASETO payload body
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`

	// For JWT only
	*jwt.RegisteredClaims
}

// NewPayload is contructor of Payload struct
func NewPayload(username string, duration time.Duration, authType string) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	IssuedAt := time.Now()
	ExpiredAt := IssuedAt.Add(duration)
	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  IssuedAt,
		ExpiredAt: ExpiredAt,
	}

	if authType == "JWT" {
		registeredClaims := jwt.RegisteredClaims{
			Issuer:    "simple-bank",
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(IssuedAt),
			ExpiresAt: jwt.NewNumericDate(ExpiredAt),
		}
		payload.RegisteredClaims = &registeredClaims
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
