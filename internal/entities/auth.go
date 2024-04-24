package entities

import "github.com/gofrs/uuid"

// Send pass to midleware.
type CtxPassKey struct{}

// Send values through middleware in context.
// TODO - move to middlewares, solve cycle import problem
type AuthContext struct {
	userID       uuid.UUID
	isRegistered bool
}

func NewAuthContext(userID uuid.UUID, isRegistered bool) AuthContext {
	return AuthContext{userID: userID, isRegistered: isRegistered}
}

func (c AuthContext) GetUserID() uuid.UUID {
	return c.userID
}

func (c AuthContext) IsRegistered() bool {
	return c.isRegistered
}
