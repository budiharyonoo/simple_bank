package token

import (
	"fmt"
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker() (Maker, error) {
	config, err := utils.LoadConfig("..")
	if err != nil {
		return nil, fmt.Errorf("minimum of JWT Secret must be %s", err.Error())
	}

	if len(config.JWTSecret) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("minimum of chacha20poly1305 secret must be %d", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{paseto.NewV2(), []byte(config.JWTSecret)}, nil
}

func (m PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration, "PASETO")
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}

func (m PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
