package token

import "time"

// token maneger
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
