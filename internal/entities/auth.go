package entities

import "github.com/gofrs/uuid"

// send pass to midleware.
type CtxPassKey struct{}

// Send values through middleware in context.

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
