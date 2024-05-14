package entities

import "github.com/gofrs/uuid"

// Send pass to midleware.
type CtxPassKey struct{}

// Send values through middleware in context.
type AuthContext struct {
	jwt          string
	isRegistered bool
	userID       uuid.UUID
}

func NewAuthContext(userID uuid.UUID, jwt string, isRegistered bool) AuthContext {
	return AuthContext{userID: userID, jwt: jwt, isRegistered: isRegistered}
}

func (c AuthContext) GetUserID() uuid.UUID {
	return c.userID
}
func (c AuthContext) GetUserJWT() string {
	return c.jwt
}

func (c AuthContext) IsRegistered() bool {
	return c.isRegistered
}
