package entities

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

// JWT for JWT token.
type JWT struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}
