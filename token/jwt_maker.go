package token

import (
	"fmt"
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecretKeyLength = 32

type JWTMaker struct {
	SecretKey string
}

func NewJWTMaker() (Maker, error) {
	config, err := utils.LoadConfig("..")
	if err != nil {
		return nil, fmt.Errorf("minimum of JWT Secret must be %s", err.Error())
	}

	if len(config.JWTSecret) < minSecretKeyLength {
		return nil, fmt.Errorf("minimum of JWT Secret must be %d", minSecretKeyLength)
	}

	return &JWTMaker{config.JWTSecret}, nil
}

func (j JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration, "JWT")
	if err != nil {
		return "", err
	}

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (j JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
