package token

import "time"

// Maker is interface for CreateToken & VerifyToken for JWT & PASETO
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
